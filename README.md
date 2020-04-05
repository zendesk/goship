# goship [![Build Status](https://travis-ci.com/zendesk/goship.svg?branch=master)](https://travis-ci.com/zendesk/goship)

Helps find and connect to particular cloud resources using defined providers.

## Usage

#### 1. Download latest version of goship from release page

[Releases](https://github.com/zendesk/goship/releases/latest)

#### 2. Move binary to /usr/local/bin directory

```mv <goship_binary> /usr/local/bin/goship && chmod +x /usr/local/bin/goship```

#### 3. Configure basic settings

```goship configure```

#### 4. Configure providers

see `config.yaml.example` for the proper configuration format

#### 5. Read help & Enjoy

```goship help```

## Examples

##### Log into instance

```
$ goship ssh kafka production
1. kafka-03c2b7c                               production   default    us-east-1d    172.18.214.1
2. kafka-04b3989                               production   default    us-east-1f    172.18.237.2
3. kafka-06b99eb                               production   default    us-east-1e    172.18.216.3
4. kafka-0b5875d                               production   default    us-east-1a    172.18.196.4
Choose your ship: 2
Logging into kafka-04b3989 (production)
```

##### Find all instances with some attribute (tag, ami, instance_type, etc.)

```
$ goship find kafka production
kafka-06b99
  ami_id              ami-04684244
  az                  us-east-1e
  dns_name
  id                  i-06b99eb666ddeeeee
  instance_type       m5.2xlarge
  key_name            infrastructure
  private_dns_name    ip-172-18-216-1.ec2.internal
  private_ip          172.18.216.1
  public_ip           <nil>
  tags
    Name              kafka-06b99
    project           default
    hostgroup         kafka
    owner             ops
    environment        production


kafka-0b587
  ami_id              ami-04684244
  az                  us-east-1a
  dns_name
  id                  i-0b5875d03d1b3caaa
  instance_type       m5.2xlarge
  key_name            infrastructure
  private_dns_name    ip-172-18-196-3.ec2.internal
  private_ip          172.18.196.3
  public_ip           <nil>
  tags
    Name              kafka-0b587
    project           default
    hostgroup         kafka
    owner             ops
    environment        production

[...]
```


##### Scp file to host by name

```
goship scp ~/test-file kafka:~/ --username ubuntu
1. kafka-03c2b7c                               production   default    us-east-1d    172.18.214.1
2. kafka-04b3989                               production   default    us-east-1f    172.18.237.2
3. kafka-06b99eb                               production   default    us-east-1e    172.18.216.3
4. kafka-0b5875d                               production   default    us-east-1a    172.18.196.4
Choose your ship: 2
Copying to kafka-04b3989 (production)
```

### Removing cache

In order to uncache existing files, just use `--uncache` flag.

## Contributors

* Marcin Matłaszek ([emate](https://github.com/emate))
* Szymon Władyka ([swladyka](https://github.com/swladyka))
* Karol Gil ([karolgil](https://github.com/karolgil))
* Leszek Charkiewicz ([lcharkiewicz](https://github.com/lcharkiewicz))

## License

Copyright 2018 Zendesk

Licensed under the [Apache License, Version 2.0](LICENSE)
