# vacuum

Deep clean your AWS account of unused resources to save you mon
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
## What is it?

Even when using infrastructure as code tools like Terraform it is easy to leave behind relic resources such as EC2 EBS volumes and ENIs.

These relics take money out of your pocket and put it in Jeff Bezos's,  Enter vacuum, to deep clean your account!

## All

To run all the commands and thoroughly deep clean your account run `vacuum all`

By default, regions `eu-west-1` and `eu-west-2` will be vacuumed.  You can override this using the regions flag:

```
vacuum all -r "us-east-1,us-east-2"
```

## Volumes

Clean up available EC2 volumes using the `volumes` command.  The volumes are not attached to anything and are just sitting there lining Jeff Bezos's pocket.  No one wants that so clean them up using:

```
vacuum volumes
```

By default, regions `eu-west-1` and `eu-west-2` will be vacuumed.  You can override this using the regions flag:

```
vacuum volumes -r "us-east-1,us-east-2"
```


## ENIs

Clean up available ENIs using the `enis` command.  The ENIs are not attached to anything and are just sitting there costing you money, clean them up using:

```
vacuum enis
```

By default, regions `eu-west-1` and `eu-west-2` will be vacuumed.  You can override this using the regions flag:

```
vacuum enis -r "us-east-1,us-east-2"
```

