# Optimizing Emergency Call Routing Using KNN and Graph Theory

## Project Overview

This project aims to improve the efficiency of emergency call routing in Saudi Arabia's telecommunication network by leveraging Linear Algebra concepts, specifically the K-Nearest Neighbors (KNN) algorithm, and graph theory. The proposed system optimizes the process of directing emergency calls from telecom sites to the nearest emergency centers, ensuring timely and efficient responses.

## Network Infrastructure

The network infrastructure for the project is modeled as follows:

- **Telecom Sites**: Over 15,000 telecom sites distributed across Saudi Arabia.
- **Hub Sites**: More than 5,000 hub sites act as intermediate nodes in the network.
- **Emergency Centers**: Over 100 emergency centers provide the endpoints for routing calls.

### Assumptions:

- Each telecom site is connected to its nearest hub site.
- Each hub site is connected to its two nearest hub sites.
- Each emergency center is connected to its five nearest hub sites.

## Objectives

1. **Implement KNN Algorithm**:

   - Identify the nearest emergency center for any given telecom site in the network.

2. **Graph Theory Application**:

   - Compute the most efficient path for calls to travel through the network to reach the correct emergency center.

## Features

### Key Functionalities

- **KNN-Based Emergency Center Identification**:

  - Given a telecom site, determine the nearest emergency center based on geographic coordinates.

- **Graph-Based Path Optimization**:

  - Use Dijkstra's algorithm to compute the shortest path for emergency call routing.

### Data Structures

- **Point**:

  - Represents a 2D coordinate with `X` and `Y` values.

- **TelecomSite, Hub, EmergencyCenter**:

  - Structures representing the nodes in the network.

- **Graph**:

  - A graph data structure that maintains nodes and edges to model network connections.

### Algorithms

1. **KNN Algorithm**:

   - Calculate distances between telecom sites and emergency centers to identify the closest center.

2. **Dijkstra's Algorithm**:

   - Determine the shortest path between nodes in the graph, ensuring efficient routing of emergency calls.

## Implementation Details

### Infrastructure Connections

- **Telecom Sites to Hubs**:

  - Each telecom site connects to its nearest hub.

- **Hub to Hub**:

  - Each hub connects to its two nearest hubs.

- **Emergency Centers to Hubs**:

  - Each emergency center connects to its five nearest hubs.

### Simulation Setup

- Data for telecom sites, hubs, and emergency centers is simulated based on hypothetical coordinates.
- The system constructs a graph representing the entire network with all connections.

## Future Enhancements

- **Scalability**:

  - Extend the network to support additional sites and centers.

- **Dynamic Routing**:

  - Implement real-time updates for network changes and traffic conditions.

- **Machine Learning**:

  - Explore predictive models for emergency call volumes to preemptively optimize routing paths.

## License

This project is licensed under the MIT License.

---

For any inquiries, please open an issue in this repository.

