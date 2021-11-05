# vacuum
```
                            ▒▒▒▒▒▒▒▒▒▒▒▒
                          ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
        ░░░░░░░░░░        ▒▒▒▒        ▒▒▒▒▒▒
      ░░░░░░░░░░░░░░      ▒▒▒▒          ▒▒▒▒
      ▒▒▒▒▒▒▒▒▒▒▒▒▒▒      ▒▒▒▒▒▒        ▒▒▒▒
  ░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░    ▓▓▒▒▒▒      ▒▒▒▒
░░░░░░░░▒▒▒▒▒▒▒▒▒▒░░░░░░░░    ▒▒▒▒      ▒▒▒▒
▒▒░░░░░░░░░░░░░░░░░░░░░░▒▒    ▒▒▒▒        ▒▒▒▒
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒        ▒▒▒▒
░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░    ▒▒▒▒        ▒▒▒▒
░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░  ▒▒▒▒▒▒        ▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒▒          ▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒            ▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░░░                ▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▓▓
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▒▒
▒▒░░░░░░░░░░░░░░░░░░░░░░▒▒                  ▒▒▒▒
▒▒▒▒▒▒▒▒░░░░░░░░░░▒▒▒▒▒▒▒▒                  ▒▒▒▒
██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██                ▒▒▒▒▒▒▒▒
████    ▒▒▒▒▒▒▒▒▒▒    ████              ▒▒▒▒▒▒▒▒▒▒▒▒
          ████                            ░░░░░░░░░░░░░░░░░░
                                        ░░░░░░░░░░░░░░░░░░
                                            ░░░░░░░░░░
```

Even when using infrastructure as code tools like Terraform it is easy to leave behind relic resources such as EC2 EBS volumes and ENIs.

These relics cost you money!  Enter vacuum, to deep clean your account!

## Volumes

Clean up available EC2 volumes using the `volumes` command.  The volumes are not attached to anything and are just sitting there lining Jeff Bezos's pocket.  No one wants that so clean them up using:

```
vacuum volumes
```

By default, regions `eu-west-1` and `eu-west-2` will be vacuumed.  You can override this using the regions flag:

```
vacuum volumes -r "us-east-1,us-east-2"
```

