
<p align="center">

![](./docs/demo.gif)
</p>

**gocps** is a simple tool which aims to help HAM radio operators creating
CPS's for their DMR radios.

## Usage / Examples

 - List all German users

   ```./cps -c DEU```

 - Show just those frome 7th region:

   ```./cps -c DEU -s "^.+7"```

 - ...And only those of name *Harald*:

   ```./cps -c DEU -s "^.+7" -n Harald```

 - Tweak the output format a bit:

   ```./cps -c DEU -s "^.+7" -n Harald -f "{{.id}},{{.call}},{{.name}}"```