## mole domains add proxy

Add a reverse proxy for a domain

### Synopsis

This command creates a reverse proxy configuration in Caddy for the specified project.
	if an empty on 0 port flag is set, MOLE_PORT_APP env variable will be used instead.

```
mole domains add proxy [project name/id] [flags]
```

### Options

```
  -d, --domain string   Domain *required
  -h, --help            help for proxy
  -p, --port int        Port *required
```

### SEE ALSO

* [mole domains add](mole_domains_add.md)	 - Add a new domain to the Caddy configuration

###### Auto generated by spf13/cobra on 26-Nov-2024
