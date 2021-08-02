import fileinput,sys,argparse

parser=argparse.ArgumentParser()

parser.add_argument('--file', help='Filename', required=True)
parser.add_argument('--key', help='The key to update', required=True)
parser.add_argument('--value', help='New value as is. Add quotes if required', required=True)

args=parser.parse_args()
tag=args.key
newValue = args.value

print ("{} {} ".format(tag, newValue))

counter = 0
for  line in fileinput.FileInput(args.file, inplace=1):
    if line.lstrip().startswith("- "+tag) or line.lstrip().startswith(tag):
        if line.index(':') <= 0:
            sys.exit("Can't find a colon in a line: " + line)
        print("{}: \"{}\" ".format(line[:line.index(':')],newValue))
        counter += 1
    else:        
        sys.stdout.write(line)
if counter < 0:
    print ("Key not found")
    sys.exit(1)
print ("Lines updated:  {}".format(counter))