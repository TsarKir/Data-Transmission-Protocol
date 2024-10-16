# Data Transmission Protocol

## Go Team 00

### My Educational Project on Golang

### Task: Implementation of a Data Transmission Protocol

### Task Description
Implement a data transmission protocol using gRPC. Each transmitted data packet should contain three fields: `session_id` (string), `frequency` (floating-point number), and the current timestamp in UTC format.

### Implementation Steps

1. **Create a .proto file**:
   Define the data schema in a `.proto` file that describes the message structure. Example:
   ```protobuf
   syntax = "proto3";

   message FrequencyData {
       string session_id = 1;
       double frequency = 2;
       string timestamp = 3; // Timestamp in UTC format
   }

   service FrequencyService {
       rpc StreamFrequencies(stream FrequencyData) returns (stream FrequencyData);
   }


2. **Generate code**:
   Use protoc to generate Go code from the .proto file. This will create the necessary structures and methods for working with gRPC.

3. **Implement the server**:
   - For each new client connection, generate a unique session_id and random values for the mean and standard deviation.
   - Log these values (to stdout or a file).
   - Send a stream of messages where frequency is randomly selected from a normal distribution defined by the generated mean and standard deviation.

4. **Implement the client**:
   - The client should receive a stream of values and gradually update its estimates for the mean and standard deviation based on the incoming data.
   - After processing a certain number of values (e.g., 50-100), the client should automatically switch to anomaly detection mode.
   - This mode should consider an anomaly coefficient passed via command line.

5. **Anomaly detection**:
   - A frequency value is considered anomalous if it deviates from the expected value by more than k × STD (where k is the anomaly coefficient).
   - All detected anomalies should be logged.

6. **Store anomalies in a database**:
   - Use an ORM to record anomalies in PostgreSQL. Define a record structure (e.g., session_id, frequency, and timestamp) and use ORM to map these fields to database columns.
