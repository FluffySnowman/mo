# Mo - S3 compatible CLI manager for object storage

Alternative to s3cmd, aws cli, etc.

Don't expect updated or correct documentation. This repository is a tool I built
which I use myself and the code just happens to be publicly available.

# Thing

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


