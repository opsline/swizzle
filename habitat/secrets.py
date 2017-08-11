#!/usr/bin/python
import sys
import boto3
import urllib2
import argparse
import traceback

parser = argparse.ArgumentParser()
parser.add_argument('-A', '--action',
        dest='action', action='store', required=False,
        choices=['get', 'create', 'delete'], default='get',
        help='Action')
parser.add_argument('-e', '--environment',
        dest='environment', action='store', required=True,
        help='Environment name')
parser.add_argument('-a', '--application',
        dest='application', action='store', required=True,
        help='Application name')
parser.add_argument('-s', '--secret-name',
        dest='secret_name', action='store', required=False,
        help='Secret name')
parser.add_argument('-v', '--secret-value',
        dest='secret_value', action='store', required=False,
        help='Secret value')
parser.add_argument('-i', '--from-file',
        dest='from_file', action='store', required=False,
        default='',
        help='Input file')
parser.add_argument('-r', '--aws-region',
        dest='aws_region', action='store', required=False,
        help='AWS region name')
parser.add_argument('-p', '--aws-profile',
        dest='aws_profile', action='store', required=False,
        help='AWS CLI profile name')
parser.add_argument('-k', '--kms-key',
        dest='ksm_key', action='store', required=False,
        default='alias/aws/ssm',
        help='KMS key alias')
parser.add_argument('-f', '--format',
        dest='format', action='store', required=False,
        choices=['env', 'json', 'yaml'], default='env',
        help='Output format')
args = parser.parse_args()

def ssm_connection():
    if args.aws_region is None:
        u = 'http://169.254.169.254/latest/meta-data/placement/availability-zone'
        try:
            region = urllib2.urlopen(u, timeout=3).read()[:-1]
        except:
            sys.stderr.write('could not read ec2 metadata - please provide AWS region name\n')
            sys.exit(1)
    else:
        region = args.aws_region
    session = boto3.Session(profile_name=args.aws_profile)
    client = session.client('ssm', region_name=region)
    return client

def get_parameter_names(client):
    param_names = []
    next_token = None
    have_more = True
    param_filters = [
        {
            'Key': 'Name',
            'Option': 'BeginsWith',
            'Values': [get_param_name_prefix()]
        }
    ]
    try:
        while have_more:
            if next_token is not None:
                response = client.describe_parameters(
                        ParameterFilters=param_filters, NextToken=next_token)
            else:
                response = client.describe_parameters(
                        ParameterFilters=param_filters)
            for param in response['Parameters']:
                param_names.append(param['Name'])
            if response.has_key('NextToken'):
                next_token = response['NextToken']
            else:
                have_more = False
        return param_names
    except:
        sys.stderr.write('ERROR: could not retrieve parameter names\n')
        sys.stderr.write(traceback.format_exc())
        sys.exit(1)

def split_array(arr, size):
     arrs = []
     while len(arr) > size:
         pice = arr[:size]
         arrs.append(pice)
         arr = arr[size:]
     arrs.append(arr)
     return arrs

def get_parameters(client, names):
    params = []
    for smaller_names in split_array(names, 10):
        try:
            response = client.get_parameters(Names=smaller_names, WithDecryption=True)
        except:
            sys.stderr.write('ERROR: could not retrieve parameters\n')
            sys.stderr.write(traceback.format_exc())
            sys.exit(1)
        for p in response['Parameters']:
            params.append(p)
    return params


def get_param_name_prefix():
    return '%s.%s.' % (args.environment, args.application)

def get_full_param_name(name):
    return '%s.%s.%s' % (args.environment, args.application, name)

def create_parameter(client):
    if args.from_file != '':
        input_file = open(args.from_file, 'r')
        for line in input_file.readlines():
            if line.find('=') > 0:
                next
            (n, v) = line.split('=')
            n = n.replace('export ', '')
            v = v[:-1]
            if v.startswith(('"', "'")) and v.endswith(('"', "'")):
                v = v[1:-1]
            fpn = get_full_param_name(n)
            sys.stderr.write('creating: %s\n' % fpn)
            try:
                client.put_parameter(
                        Name=fpn, Value=v, Description=n,
                        Type='SecureString', KeyId=args.ksm_key, Overwrite=True)
            except:
                sys.stderr.write('ERROR: could not create parameter\n')
                sys.stderr.write(traceback.format_exc())
                sys.exit(1)
    else:
        fpn = get_full_param_name(args.secret_name)
        sys.stderr.write('creating: %s\n' % fpn)
        try:
            client.put_parameter(
                    Name=fpn, Value=args.secret_value,
                    Type='SecureString', KeyId=args.ksm_key, Overwrite=True)
        except:
            sys.stderr.write('ERROR: could not create parameter\n')
            sys.stderr.write(traceback.format_exc())
            sys.exit(1)

def delete_parameter(client):
    fpn = get_full_param_name(n)
    sys.stderr.write('deleting: %s\n' % fpn)
    try:
        client.delete_parameter(Name=fpn)
    except:
        sys.stderr.write('ERROR: could not delete parameter\n')
        sys.stderr.write(traceback.format_exc())
        sys.exit(1)

def main():
    client = ssm_connection()
    if args.ksm_key == '':
        args.ksm_key = 'alias/secrets/%s' % args.environment
    if args.action == 'get':
        if args.secret_name is None:
            names = get_parameter_names(client)
        else:
            names = [get_full_param_name(args.secret_name)]
        params = get_parameters(client, names)
        prefix = get_param_name_prefix()
        p = [
                {
                    'name': str(param['Name']).replace(prefix, ''),
                    'value': str(param['Value'])
                } for param in params
            ]
        if args.format == 'env':
            for param in p:
                print("export %s='%s'" % (param['name'], param['value']))
        elif args.format == 'json':
            import json
            print(json.dumps(p, indent=2, sort_keys=True))
        elif args.format == 'yaml':
            import yaml
            print(yaml.dump(p, indent=2, default_flow_style=False))

    elif args.action == 'create':
        create_parameter(client)
    elif args.action == 'delete':
        delete_parameter(client)

if __name__ == "__main__":
    main()
