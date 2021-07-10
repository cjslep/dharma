# dharma

> Ode To Eve, circa Scarcity-Era, Sovicou
>
> Use Dharma and then go undock,
> wreck ships, mine some ore, make small talk
> at the end of the day
> see how much you've been paid
> still less than male Dancers' stock

Dharma is the premier community-focused software for independent EvE Online
corporations, integrating with ESI for additional functionality. It aims to
be a first-class tool around community building, event-planning, asset
management, intel-sharing, KoS coordination, and industry/logistics planning
platform.

**Dharma is pre-alpha software still under development, not yet ready for use.
Due to its ambitious nature, it will take some time. This is a garden being
tended to, not a quarry being excavated.**

## Federation In A Nutshell

Dharma is built using the ActivityPub federation protocol. This means that each
corporation has the option to granularly control the data it shares with other
corporations also running Dharma. Corporations are more freely able to associate
and disassociate with each other in this way -- including, but not limited to,
Alliances and coalitions -- while still preserving their own unique sense of
community. This allows corporations using Dharma to have tools that readily
serve a solo player corp -- lowering the barrier for to communicate with other
small corps -- yet scale in effectiveness to a federating size that rivals
a traditional centralized Alliance.

That means one-capsuleer or few-capsuleer corporations have new channels through
which they can discover other micro-corporations, and lower the barrier towards
building smaller and perhaps more chaotic political networks of mutual aid.
Large corporations are also welcome to adopt the software and more readily
interact with smaller ones.

This software only lowers the barrier to change and keeps the bonds between
groups lubricated, yet still it is the capsuleer that makes the corporation.

## Requirements

Installation is somewhat technical:

- A domain name.
- A server machine.
- A postgres database.
- An Omega account in EvE Online, required to obtain your ESI key. You need your
  own ESI key, so that CCP Games has granular control in dealing with others who
  would otherwise abuse the ESI API, allowing you to remain unaffected by
  others' actions.

Instructions on the above will be elaborated upon, later.

## Installation

A release is not yet available. Do not install this software at this time. Once
a release is available, these instructions will be updated.

## Features

None, this is still under development. The goal is to support:

| Feature                   | Status |
| ------------------------- | ------ |
| Core Account Management   | 🔨     |
| Local Forum               | 🔨     |
| Federating Forum Messages | 🔨     |
| Corporation Standings     | 🕖     |
| KoS & Justice Management  | 🕖     |
| Federating KoS & Justice  | 🕖     |
| Federation Controls       | 🕖     |
| Federation Audit Log      | 🕖     |
| Calendar Integration      | 🕖     |
| Federating Calendar       | 🕖     |
| Intel (3rd party?)        | 🕖     |
| Federating Intel          | 🕖     |
| Fitting (3rd party?)      | 🕖     |
| Federating Fittings       | 🕖     |
| Asset ESI                 | 🕖     |
| Industry & Logistics      | 🕖     |
| Federating Indy & Logi    | 🕖     |

For some of these, I would like to look into interfacing with well-established
tools for integration (ex: for fitting and intel-sharing). This software will
also include focusing on tooling for NRDS RoE and lore/roleplay communities, but
more research is needed to understand feature requirements.

## Feature Requests, Reporting Bugs, Contributing Code

If you wish to discuss this software and/or its features, please see the
[`CONTRIBUTING.md`](./CONTRIBUTING.md) file.
