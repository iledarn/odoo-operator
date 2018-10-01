package cluster

const odooBasicFmt = `
[options]

data_dir = %s

without_demo = %s
server_wide_modules = %s

proxy_mode = True

db_name = %s
db_template = %s

list_db = %s
dbfilter = %s

unaccent = True
publisher_warranty_url = %s

log_handler = %s

email_from = %s
smtp_server = %s
smtp_port = %s
smtp_ssl = %s
smtp_user = %s
smtp_password = %s

`

const odooCustomSection = `
[backups]
backupfolder = %s

[integrator]
integrator_warranty_url = %s

`

const odooPsqlSecretFmt = `
# hostname:port:database:username:password
# If an entry needs to contain : or \, escape this character with \
%s:%s:*:%s:%s
`
const odooAdminSecretFmt = `%s

`
