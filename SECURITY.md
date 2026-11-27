# Security policy

## Supported versions

Tokenkitty is pre-alpha. Until 1.0, only the latest release gets security fixes.

## Reporting a vulnerability

**Please don't open a public issue.**

Report it privately through GitHub: go to the [Security tab](https://github.com/tokenkitty/tokenkitty/security) and choose **Report a vulnerability**. Include:

- what is affected (collector, server or dashboard) and which version
- steps to reproduce, or a proof of concept
- the impact you expect

You'll get an acknowledgement within 3 working days. We'll keep you updated as we fix it, and credit you in the advisory unless you'd rather stay anonymous.

## Scope

Especially important to us:

- any way the collector could read or send prompt text, model output, file contents, file paths or commands
- a device token being used to submit events as another device or person
- reading another person's data with the `self` role
- leaking integration credentials stored by the server
