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
  - [kubelet](#kubelet)
  - [kube-controller-manager](#kube-controller-manager)
  - [cloud-controller-manager](#cloud-controller-manager)
  - [kubeadm](#kubeadm)
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
your Kubernetes clusters will lack features that require integration with the underlying infrastructure/cloud provider.

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
