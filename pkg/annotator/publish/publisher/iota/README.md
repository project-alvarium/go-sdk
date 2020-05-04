# IOTA Tangle Publisher

## Table of Contents
- [Quick Start](#quick-start)
  - [Non-Linux Hosts](#non-linux-hosts)
  - [Setup Guide](#setup-guide)
- [IOTA](#iota)
  - [Overview](#overview)
  - [IOTA Networks](#iota-networks)
    - [Terminology](#terminology)
      - [Nodes](#nodes)
      - [Seeds](#seeds)
      - [Addresses](#addresses)
      - [IOTA Tokens](#iota-tokens)
      - [Transactions](#transactions)
      - [Transfer Bundles](#transfer-bundles)
      - [IOTA's Glossary for Other Related Terms](#iotas-glossary-for-other-related-terms)
    - [Public](#public)
      - [Mainnet](#mainnet)
      - [Devnet](#devnet)
    - [Private](#private)
- [SDK](sdk)
  - [Seed](#seed)
  - [Depth](#depth)
  - [Minimum Weight Magnitude](#minimum-weight-magnitude)
  - [SDK Contract](#sdk-contract)

## Quick Start

For testing, you may not want to publish test data on one of the public networks. `one-command-tangle` serves as a one-stop-shop where you can run your own local private IOTA network using a single Docker command. This application is useful when testing applications without risking any monetary value.

### Non-Linux Hosts

By default, `one-command-tangle` uses a Docker `host` network. The `host` networking driver [only works on Linux hosts](https://docs.docker.com/network/host/), and is not supported on Docker Desktop for Mac, Docker Desktop for Windows, or Docker EE for Windows Server.

To get `one-command-tangle` running on a non-Linux host, you will need to make the following changes to the `one-command-tangle` docker-compose.yml file.

- Replace `network_mode` value `host` in [compass service's network mode](https://github.com/iota-community/one-command-tangle/blob/00c5cc70c5417b26738c954bf0928fc703dfb4ab/docker-compose.yml#L5) with `bridge`.
- Replace `network_mode` value `host` in [IRI service's network mode](https://github.com/iota-community/one-command-tangle/blob/00c5cc70c5417b26738c954bf0928fc703dfb4ab/docker-compose.yml#L24) with `bridge`. 

### Setup Guide

For more information on using `one-command-tangle`, refer to this [one-command-tangle setup guide](https://docs.iota.org/docs/utils/0.1/community/one-command-tangle/overview).

## IOTA

IOTA is an open-source distributed ledger technology that allows connected devices to transfer data and [IOTA tokens](#iota-token) among each other for zero fees.
 
### Overview

Setting up an IOTA Network for your use-case may vary if complexity. If you are unfamiliar with IOTA Networks, checking out [IOTA overview](https://docs.iota.org/docs/getting-started/0.1/introduction/overview) is recommended. 
 
### IOTA Networks

#### Terminology

Prior to configuring an IOTA Network, it is important to understand some of the base entities which are used in
 Networks of all variations.
 
##### Nodes

Nodes are the core of an IOTA network. They run the [node software](https://docs.iota.org/docs/node-software/0.1/introduction/overview) that gives them read and write access to the Tangle. Like any distributed system, nodes are connected to others called neighbors to form an IOTA network. When one node, no matter where it is in the world, receives a [transaction](#transaction), it will try to forward it to all of its neighbors. This way, all nodes eventually validate all transactions and store them in the local copy of the Tangle called a ledger.

For more information on nodes, refer to the [nodes documentation](https://docs.iota.org/docs/getting-started/0.1/network/nodes).
 
##### Seeds

A seed is a unique password that gives you the ability to prove your ownership of either messages and/or any [IOTA tokens](#iota-token) that help on your [addresses](#address).
 
For more information on seeds, refer to the [seeds documentation](https://docs.iota.org/docs/getting-started/0.1/clients/seeds).

##### Addresses

An address is like an account that belongs to a seed and that has a 0 or greater balance of [IOTA Tokens](#iota-token). Addresses are the public half of a public/private key pair. To transfer [IOTA Tokens](#iota-token) from one address to another, you sign a transaction with the private key to prove to nodes that you own it. As such you can share addresses with anyone because eon the seed owner knows the private key.
  
For more information on addresses, refer to the [address documentation](https://docs.iota.org/docs/getting-started/0.1/clients/addresses).

##### IOTA Tokens

An IOTA token is a unit of value (i, Ki, Mi, Gi, Ti, Pi) that can be transferred over an IOTA network in transfer bundles.
 
For more information on IOTA Tokens, refer to the [IOTA tokens documentation](https://docs.iota.org/docs/getting-started/0.1/clients/token).

##### Transactions

A transaction is a single transfer instruction that can either withdraw [IOTA tokens](#iota-token) from an [address](#address), deposit them into an address, or have zero-value (contain data, a message, or a signature). If you want to send anything to an IOTA network, you must send it to a [node](#node) as a transaction.

For more information on transactions, refer to the [transactions documentation](https://docs.iota.org/docs/getting-started/0.1/transactions/transactions).

##### Transfer Bundles

A bundle is a group of transactions that rely on each other's validity. For example, a transaction that deposits [IOTA tokens](#iota-token) into an address relies on another transaction to withdraw those IOTA tokens from another address. Therefore, those transactions must be in the same bundle.

For more information on bundles, refer to the [bundles documentation](https://docs.iota.org/docs/getting-started/0.1/transactions/bundles).

##### IOTA's Glossary for Other Related Terms

IOTA Networks contain a lot of terms/concepts that are used throughout networks of all sorts. IOTA has a glossary which can be accessed [here](https://docs.iota.org/docs/getting-started/0.1/references/glossary).

#### Public

IOTA has two public networks of nodes, and each one has its own Tangle to which nodes can attach transactions.

##### Mainnet

IOTA network that uses IOTA Tokens, which are traded on cryptocurrency exchanges.

For more information on public Mainnet IOTA Networks, refer to [IOTA Networks Official Documentation](https://docs.iota.org/docs/getting-started/0.1/network/iota-networks#mainnet)

##### Devnet

Devnet is similar to the Mainnet, except the tokens are free and it takes less time and computational power to create and send a transaction.
 
On this network, you can test your applications and build proofs of concept that use free Devnet tokens.
 
For more information on public Devnet IOTA Networks, refer to [IOTA Networks Official Documentation](https://docs.iota.org/docs/getting-started/0.1/network/iota-networks#devnet).

#### Private

A private Tangle is an IOTA network that you control and that contains only nodes that you know. A private Tangle uses the same technology as the public IOTA networks, except you control it by running an open-source implementation of the Coordinator called Compass. You can use Compass to allow nodes to reach a consensus on transactions attached to your private Tangle. If Compass stops, no transactions in your IOTA network will be confirmed until it starts again.
    
For more information on setting up a private Tangle, refer to this [set up a private Tangle guide](https://docs.iota.org/docs/compass/0.1/how-to-guides/set-up-a-private-tangle).

## SDK

### Seed

A seed is a unique password that gives you the ability to prove ownership of either messages and/or any [IOTA tokens](https://docs.iota.org/docs/getting-started/0.1/clients/token) that are held on you [addresses](https://docs.iota.org/docs/getting-started/0.1/clients/addresses).

[Official seed documentation](https://docs.iota.org/docs/getting-started/0.1/clients/seeds)

### Depth

When sending a transaction to a [node](https://docs.iota.org/docs/getting-started/0.1/network/nodes), you can specify a depth argument, which defines how many milestones in the past the node starts the [tip selection algorithm](https://docs.iota.org/docs/getting-started/0.1/network/the-tangle#tip-selection). The greater the depth, the farther back in the [Tangle](https://docs.iota.org/docs/getting-started/0.1/network/the-tangle) the node starts. A greater depth increases the time that node takes to complete tip selection, making them use more computational power.

[Official depth documentation](https://docs.iota.org/docs/getting-started/0.1/transactions/depth)

### Minimum Weight Magnitude

When doing [proof-of-work](https://docs.iota.org/docs/getting-started/0.1/transactions/proof-of-work) for a [transaction](https://docs.iota.org/docs/getting-started/0.1/transactions/transactions), you can specify a minimum weight magnitude (MWM), which defines how much work is done. Nodes accept a certain minimum weight magnitude, depending on their configuration. If you use a minimum weight magnitude that is too low, nodes will reject the transaction as invalid.

[Official minimum weight magnitude documentation](https://docs.iota.org/docs/getting-started/0.1/network/minimum-weight-magnitude)

### SDK Contract

The SDK contains an [IOTA Tangle publisher contract](sdk/contract.go). This contract is an abstraction which facilitates unit testing such that methods which require interaction over HTTP with an IOTA Tangle, can be stubbed. 

The SDK provides a [default implementation](sdk/iota/iota.go). This default implementation leverages the underlying IOTA ledger dependency which serves as a client to an IOTA Tangle. 
