# ngfshare

Simple file sharing server utilizing nginx X-Accel file serving.

## Setup

Example config.json
```json
{
    "Port": 8987,
    "Address": "127.0.0.1",
    "DBpath": "/var/lib/ngfshare/database.db",
    "StoreDir": "/var/lib/ngfshare/store",
    "HTMLTemplateDir": "/var/lib/ngfshare/templates",
    "UrlPrefix": "https://share.yourdomain.xyz",
    "IdLen": 5,
    "AuthKeyLen": 30
}
```

Example nginx server configuration
```cfg
server {
    listen 443 ssl;
    server_name share.yourdomain.xyz;
    root /nonexistent;

    location / {
        proxy_set_header Host $host; # Not actually required
        proxy_pass http://127.0.0.1:8987;
    }
    location /store/ {
        internal;
        alias /var/lib/ngfshare/store/;
    }
}
```

Caching based only on file id can be enabled with this example configuration
```cfg
location ~ ^/-([a-zA-Z0-9]+) {
    # Other caching options are omitted here...
    proxy_cache_revalidate on;
    proxy_cache_key "$host$1";
    # Other options omitted...
}
```

Then run with `./ngfshare -config /path/to/config.json`. To generate an auth key, run `./ngfshare -config /path/to/config.json -genauth`.

## Api
### Upload
Files can be uploaded by making a POST request to `/api/upload` with the `Authorization` header set to the auth key generated above.

Example
```sh
curl 'https://share.yourdomain.xyz/api/upload' -H 'Authorization: authkey' -F 'file=@img.png'
```

This returns a json dict
```json
{
  "id": "A5HjC",
  "filename": "img.png",
  "url": "https://share.yourdomain.xyz/-A5HjC/img.png",
  "url_short": "https://share.yourdomain.xyz/-A5HjC",
  "delete_url": "https://share.yourdomain.xyz/api/delete/A5HjC"
}
```

### Delete
Make a POST request to `/api/delete/{id}` with the `Authorization` header set.

Example
```sh
curl -X POST 'https://share.yourdomain.xyz/api/delete/A5HjC' -H 'Authorization: authkey'
```

Returns
```json
{
  "status": "OK"
}
```

### Web
There is a web interface for it as well. Login with the auth key. Once logged in, a list of files are shown. New files can be uploaded, and deleted from the web interface. You can also probably tell I'm not a web designer.

## Inspiration
Inspired by https://github.com/mtlynch/picoshare. This saves the files into the sqlite3 database. I didn't like that, so thats why I made this.
