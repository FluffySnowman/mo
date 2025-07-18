# Mo - S3 compatible CLI manager for object storage

Alternative to s3cmd, aws cli, etc.

Don't expect updated or correct documentation. This repository is a tool I built
which I use myself and the code just happens to be publicly available.

Table of contents:

<!--toc:start-->
- [Mo - S3 compatible CLI manager for object storage](#mo-s3-compatible-cli-manager-for-object-storage)
- [Installation](#installation)
    - [From Release (download executable):](#from-release-download-executable)
    - [Build from source:](#build-from-source)
- [Thing for the thing of the thing's thing](#thing-for-the-thing-of-the-things-thing)
<!--toc:end-->

# Installation

### From Release (download executable):

- Curl 

```bash
sudo curl -L -o /usr/local/bin/mo https://github.com/FluffySnowman/mo/releases/download/v1.0.0/mo_x86_64-unknown-linux-gnu && sudo chmod +x /usr/local/bin/mo
```

- Wget

```bash
sudo wget -O /usr/local/bin/mo https://github.com/FluffySnowman/mo/releases/download/v1.0.0/mo_x86_64-unknown-linux-gnu && sudo chmod +x /usr/local/bin/mo
```


### Build from source:

```bash 
git clone https://github.com/fluffysnowman/mo
cd mo 
make build install
```

and the `mo` executable will be present in `/usr/local/bin/mo`. make sure this
is in your `$PATH`

# Thing for the thing of the thing's thing

```bash
mo - s3-compatible object storage manager


usage:
    mo [options] command [args...]


opts:
    -endpoint URL     s3 endpoint url
    -region REGION    s3 region (default: us-east-1)
    -key KEY          access key id
    -secret SECRET    secret access key
    -bucket BUCKET    default bucket name
    -insecure         use http instead of https
    -r                recursive copy (for cp command)


cmds:
    buckets           list all buckets
    ls [prefix]       list objects in bucket
    cp SOURCE DEST    copy/upload/download files
    rm OBJECT         remove(delet) object
    mv SOURCE DEST    move/rename object
    stat OBJECT       get object metadata
    mb [bucket]       make bucket
    rb [bucket]       remove bucket


doc/example(s(?)):
    mo -endpoint s3.amazonaws.com -key YOUARERETARDEDLOLHEH -secret wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY -bucket somebucket ls
    mo -endpoint localhost:9000 -key minioadmin -secret minioadmin -bucket test cp ./file.txt remote/file.txt
    mo buckets
    mo ls myfolder/
    mo cp ./localfile.txt remote/path/file.txt
    mo cp remote/file.txt ./localfile.txt
    mo -r cp ./mydirectory/ remote/backup/
    mo rm remote/file.txt
    mo mv old/path.txt new/path.txt
    mo stat remote/file.txt


configuration file if you can't handle typing credentials unlike gigachads:
    create a file 'mo.conf' in the $CWD/$PWD w/e 
    
    endpoint s3.amazonaws.com
    region us-east-1
    access_key AKIAIOSFODNN7EXAMPLE
    secret_key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    bucket mybucket


environment variables (overrides previouus opts):
    MO_ENDPOINT, MO_REGION, MO_ACCESS_KEY, MO_SECRET_KEY, MO_BUCKET

```


eee


