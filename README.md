# Kraken - distributed file server

Code written by lobotomized squirrel squad under tons of cocaine.

Don't use it in your real-world infrastructure, to avoid getting awkwards.

## How to
Install dependencies:
- MongoDB
- RabbitMQ

Sample configuration located under configs folder.
Default config path is: /etc/kraken/kraken.conf

### Run lonely node
Just configure that guy as you want, then build it:
```bash
make
```

And run with your custom configuration:
```bash
./kraken --config %your_config_file_path%
``` 

### Run multiply nodes
Configure like lonely node, but specify same cluster id, and then just run it on multiple
hosts.



