# goship

Helps find and connect to particular cloud resources using defined providers.

## Usage

### 1. Download latest version of goship from release page

[Releases](https://github.com/zendesk/goship/releases/latest)

### 2. Move binary to /usr/local/bin directory

```mv <goship_binary> /usr/local/bin/goship && chmod +x /usr/local/bin/goship```

### 3. Configure basic settings

```goship configure```

### 4. Configure providers

see `config.yaml.example` for the proper configuration format,

### 5. Read help & Enjoy

```goship help```

### Removing cache

In order to uncache existing files, just use `--uncache` flag.

### Contributors

* Marcin Matłaszek ([emate](https://github.com/emate))
* Szymon Władyka ([swladyka](https://github.com/swladyka))
* Karol Gil ([karolgil](https://github.com/karolgil))
* Leszek Charkiewicz ([lcharkiewicz](https://github.com/lcharkiewicz))

## License

Copyright 2018 Zendesk

Licensed under the [Apache License, Version 2.0](LICENSE)
