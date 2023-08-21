############################################
# Stage 1: Compile the binary in a containerized Golang environment
############################################
FROM golang:1.19 as build

# Copy the source files from the host
COPY . /kvstore
COPY vendor/ /kvstore/vendor/

# Set the working directory to the same place we copied the code
WORKDIR /kvstore

# Create an empty transaction log file.
RUN touch /tmp/translog.log

# Build the binary!
RUN CGO_ENABLED=0 GOOS=linux go build -a -o kvstore

############################################
# Stage 2: Build the container image proper
############################################
FROM scratch

# Copy the binary from the host.
COPY --from=build /kvstore/kvstore .

# Copy the log so that /tmp gets created (scratch does not have a /tmp neither mkdir).
COPY --from=build /tmp/translog.log tmp/

# Tell Docker to execute this command on a `docker run`.
CMD ["/kvstore"]
