# Kubernetes vSphere Cloud Provider

This is the official documentation for the Kubernetes vSphere Cloud Provider.

## Table of Contents

- [Introduction](#introduction)
- [Kubernetes Concepts](#concepts)
  - [Cloud Providers](#cloud-providers)
  - [In-Tree Cloud Providers](#in-tree-cloud-providers)
  - [Out-of-Tree Cloud Providers](#out-of-tree-cloud-providers)
- [Components and Tools](#components-and-tools)
  - [VM](#VM)
  - [vSphere](#vsphere)
  - [ESXi](#esxi)
  - [vCenter](#vcenter)
  - [govmomi](#govmomi)
  - [Kubernetes](#kubernetes)
  - [kube-apiserver](#kube-apiserver)
  - [kube-controller-manager](#kube-controller-manager)
  - [kube-scheduler](#kube-scheduler)
  - [kubelet](#kubelet)
  - [kube-proxy](#kube-proxy)
  - [cloud-controller-manager](#cloud-controller-manager)
  - [etcd](#cloud-controller-manager)
- [Architecture](#architecture)
  - [Kubernetes using the In-Tree vSphere Provider](#kubernetes-using-the-in-tree-vsphere-provider)
  - [Kubernetes using the Out-of-Tree vSphere Provider](#kubernetes-using-the-out-of-tree-vsphere-provider)
- [vSphere Integrations](#vsphere-integrations)
  - [Kubernetes Nodes](#kubernetes-nodes)
  - [Kubernetes Zones/Regions Topology](#kubernetes-zones-regions-topology)
  - [Kubernetes LoadBalancers](#kubernetes-loadbalancers)
  - [Kubnernetes Routes](#kubernetes-routes)
- [Example Configurations and Manifests](#example-manifests)
- [Tutorials](#tutorials)
- [Addons](#addons)
  - [Storage](#storage)
  - [Networking](#networking)


# Introduction

This is the official documentation for the Kubernetes vSphere cloud provider integration. This document covers key concepts,
features, known issues, installation requirements and steps for Kubernetes clusters running on vSphere. Before reading this
document, it's worth noting that a Kubernetes cluster can run on vSphere without the cloud provider integration enabled, however,
your Kubernetes clusters will not have features that require integration with the underlying infrastructure/cloud provider.

## Kubernetes Concepts

Before diving into how to operate and install your Kubernetes clusters on vSphere, it's useful to go over some well-known concepts
in Kubernetes.

### Cloud Providers

The cloud provider, from the perspective of Kubernetes, refers to the infrastructure provider for cluster resources such as
nodes, load balancers, routes, etc. Various Kubernetes components are capable of communicating with the underlying cloud provider
via APIs in order to call operations (create, update, delete) against resources required for a cluster. This capabilitiy is what
we refer to as the cloud provider integration.

As of writing this, there are two modes of cloud provider integrations: in-tree and out-of-tree. More details on these two modes below.

### In-Tree Cloud Providers

In-tree cloud providers refers to cloud provider integrations that are directly compiled and built into the core Kubernetes components.
This also means that the integration is also developed within the same git tree as Kubernetes core. As a result, updates to the cloud
provider integration must also be released at the same cadeance as the main Kubernetes release.

![In-Tree Cloud Provider Architecture](/docs/images/in-tree-arch.png "Kubernetes In-Tree Cloud Provider Architecture - from k8s.io/website")


### Out-of-Tree Cloud Providers

Out-of-tree cloud provider refers to integrations that can be developed, built and released independent of Kubernetes core. This requires
adding a new component to the cluster called the cloud-controller-manager. The cloud-controller-manager is responsible for running all the
cloud-specific control loops that were previously run in core components like the kube-controller-manager and the kubelet.


![Out-of-Tree Cloud Provider Architecture](/docs/images/out-of-tree-arch.png "Kubernetes Out-of-Tree Cloud Provider Architecture - from k8s.io/website")


### In-Tree vs Out-of-Tree

As of writing this, in-tree cloud providers are only supported for historical reasons. In the early development stages of Kubernetes, implementing
cloud providers natively (in-tree) was the most viable solution. Today, with many infrastructure providers supporting Kubernetes, new cloud providers
are requried to be out-of-tree in order to grow the project sustainably. For the existing in-tree cloud providers, there's an effort to extract/migrate
clusters to use out-of-tree cloud providers, see [this KEP](https://github.com/kubernetes/enhancements/blob/master/keps/sig-cloud-provider/20190125-removing-in-tree-providers.md) for more detais.

For Kubernetes clusters on vSphere, both in-tree and out-of-tree modes of operation are supported. However, the out-of-tree vSphere
cloud provider is strongly recommended as future releases of Kubernetes will remove support for all in-tree cloud providers.
Regardless, this document will cover both the in-tree and out-of-tree vSphere integration for Kubernetes.

## Components and Tools

Before diving into Kubernetes on vSphere, it's important to cover some key components and tools.The following section introduces key components
and tools that are part of any Kubernetes cluster running on vSphere. If you are familiar with Kubernetes and vSphere, you can skip this section.

### VM

<to-do myles>

### vSphere

<to-do myles>

### ESXi

<to-do myles>

### vCenter

<to-do myles>

### govmomi

<to-do andrew>

govmomi is a Go library for interacting with VMware vSphere APIs - ESXi and/or vCenter

### Kubernetes

Kubernetes (K8s) is an open-source system for automating deployment, scaling, and management of containerized applications.

(source: kubernetes.io)

### kube-apiserver

The Kubernetes API server validates and configures data for the api objects which include pods, services, replicationcontrollers, and others.
The API Server services REST operations and provides the frontend to the cluster’s shared state through which all other components interact.

(source: kubernetes.io)

### kube-controller-manager

The Kubernetes controller manager is a daemon that embeds the core control loops shipped with Kubernetes. In applications of robotics and automation,
a control loop is a non-terminating loop that regulates the state of the system. In Kubernetes, a controller is a control loop that watches the
shared state of the cluster through the apiserver and makes changes attempting to move the current state towards the desired state. Examples of
controllers that ship with Kubernetes today are the replication controller, endpoints controller, namespace controller, and serviceaccounts controller.

(source: kubernetes.io)

### kube-scheduler

The Kubernetes scheduler is a policy-rich, topology-aware, workload-specific function that significantly impacts availability, performance, and capacity.
The scheduler needs to take into account individual and collective resource requirements, quality of service requirements, hardware/software/policy
constraints, affinity and anti-affinity specifications, data locality, inter-workload interference, deadlines, and so on. Workload-specific requirements
will be exposed through the API as necessary.

(source: kubernetes.io)

**NOTE**: kube-scheduler will never ask the cloud provider for any information pertaining to scheduling, however, it may depend on information on resources
that were placed by other components like the kubelet.

### kubelet

The kubelet is the primary “node agent” that runs on each node. The kubelet works in terms of a PodSpec. A PodSpec is a YAML or JSON object that
describes a pod. The kubelet takes a set of PodSpecs that are provided through various mechanisms (primarily through the apiserver) and ensures
that the containers described in those PodSpecs are running and healthy. The kubelet doesn’t manage containers which were not created by Kubernetes.

(source: kubernetes.io)

### kube-proxy

The Kubernetes network proxy runs on each node. This reflects services as defined in the Kubernetes API on each node and can do simple TCP, UDP,
and SCTP stream forwarding or round robin TCP, UDP, and SCTP forwarding across a set of backends. Service cluster IPs and ports are currently found
through Docker-links-compatible environment variables specifying ports opened by the service proxy. There is an optional addon that provides cluster DNS
for these cluster IPs. The user must create a service with the apiserver API to configure the proxy.

(source: kubernetes.io)

**NOTE**: kube-proxy will never ask the cloud provider for any information pertaining to network proxy, however, it may depend on information on resources
that were placed by other components like the kube-controller-manager.

### cloud-controller-manager

The Kubernetes cloud controller manager is a daemon that embeds the cloud specific control loops shipped with Kubernetes. Each cloud provider can run
the cloud controller manager as an addon to their cluster. The cloud-controller-manager defines a specification (Go interface) that must be implemented
by every cloud provider. Various expectations and behaviors of the cluster can be tuned based on the implementation set by the cloud provider.
Cloud providers are also free to run custom controllers as part of the cloud-controller-manager.

(source: kubernetes.io)

### etcd

Consistent and highly-available key value store used as Kubernetes’ backing store for all cluster data.

(source: kubernetes.io)
