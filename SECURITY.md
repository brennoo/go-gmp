# Security Policy

go-gmp is a client library for driving Greenbone/OpenVAS vulnerability
scanners. Given what it's used for, we take security reports seriously and
would rather hear about a problem privately than have it discovered in the
wild.

## Supported versions

This project is pre-1.0: there is no long-term support branch, and only the
latest commit on the default branch (and the most recent tagged release, once
one exists) is supported with security fixes. If you're running an older
version, please upgrade before reporting an issue, if practical, so the fix
lands somewhere you can actually consume it.

## Reporting a vulnerability

Please do not open a public GitHub Issue. A public report discloses the
vulnerability to everyone running this library before a fix exists.

Instead, report it privately through
[GitHub's private security advisory feature](https://github.com/brennoo/go-gmp/security/advisories/new)
on the `brennoo/go-gmp` repository, which lets maintainers and reporters
collaborate on a fix before any public disclosure.

When reporting, include as much of the following as you can:

- The affected version or commit.
- A description of the vulnerability and its impact.
- Steps to reproduce, or a minimal proof of concept.

## What to expect

This is a community-maintained project without a dedicated security team, so
please be patient. You should expect an acknowledgement of your report,
followed by ongoing communication as the issue is investigated and, if
confirmed, fixed and released. We aim to keep you informed throughout rather
than commit to a specific response time.

## Scope

This policy covers vulnerabilities in the go-gmp library itself — for
example, a way a malicious or malformed GMP response could cause memory
corruption, wire-protocol handling that violates connection lifecycle
guarantees, or improper TLS certificate validation.

It does **not** cover vulnerabilities in Greenbone Vulnerability Manager
(`gvmd`), OpenVAS/GVM scan engines, or the GMP protocol itself. Those should
be reported directly to the [Greenbone project](https://www.greenbone.net/),
following their own disclosure process.
