// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.command.remote.outputs;

import com.pulumi.command.remote.outputs.ProxyConnection;
import com.pulumi.core.annotations.CustomType;
import java.lang.Double;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class Connection {
    /**
     * @return SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.
     * 
     */
    private @Nullable String agentSocketPath;
    /**
     * @return Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10
     * 
     */
    private @Nullable Integer dialErrorLimit;
    /**
     * @return The address of the resource to connect to.
     * 
     */
    private String host;
    /**
     * @return The password we should use for the connection.
     * 
     */
    private @Nullable String password;
    /**
     * @return The port to connect to.
     * 
     */
    private @Nullable Double port;
    /**
     * @return The contents of an SSH key to use for the connection. This takes preference over the password if provided.
     * 
     */
    private @Nullable String privateKey;
    /**
     * @return The password to use in case the private key is encrypted.
     * 
     */
    private @Nullable String privateKeyPassword;
    /**
     * @return The connection settings for the bastion/proxy host.
     * 
     */
    private @Nullable ProxyConnection proxy;
    /**
     * @return The user that we should use for the connection.
     * 
     */
    private @Nullable String user;

    private Connection() {}
    /**
     * @return SSH Agent socket path. Default to environment variable SSH_AUTH_SOCK if present.
     * 
     */
    public Optional<String> agentSocketPath() {
        return Optional.ofNullable(this.agentSocketPath);
    }
    /**
     * @return Max allowed errors on trying to dial the remote host. -1 set count to unlimited. Default value is 10
     * 
     */
    public Optional<Integer> dialErrorLimit() {
        return Optional.ofNullable(this.dialErrorLimit);
    }
    /**
     * @return The address of the resource to connect to.
     * 
     */
    public String host() {
        return this.host;
    }
    /**
     * @return The password we should use for the connection.
     * 
     */
    public Optional<String> password() {
        return Optional.ofNullable(this.password);
    }
    /**
     * @return The port to connect to.
     * 
     */
    public Optional<Double> port() {
        return Optional.ofNullable(this.port);
    }
    /**
     * @return The contents of an SSH key to use for the connection. This takes preference over the password if provided.
     * 
     */
    public Optional<String> privateKey() {
        return Optional.ofNullable(this.privateKey);
    }
    /**
     * @return The password to use in case the private key is encrypted.
     * 
     */
    public Optional<String> privateKeyPassword() {
        return Optional.ofNullable(this.privateKeyPassword);
    }
    /**
     * @return The connection settings for the bastion/proxy host.
     * 
     */
    public Optional<ProxyConnection> proxy() {
        return Optional.ofNullable(this.proxy);
    }
    /**
     * @return The user that we should use for the connection.
     * 
     */
    public Optional<String> user() {
        return Optional.ofNullable(this.user);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(Connection defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String agentSocketPath;
        private @Nullable Integer dialErrorLimit;
        private String host;
        private @Nullable String password;
        private @Nullable Double port;
        private @Nullable String privateKey;
        private @Nullable String privateKeyPassword;
        private @Nullable ProxyConnection proxy;
        private @Nullable String user;
        public Builder() {}
        public Builder(Connection defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.agentSocketPath = defaults.agentSocketPath;
    	      this.dialErrorLimit = defaults.dialErrorLimit;
    	      this.host = defaults.host;
    	      this.password = defaults.password;
    	      this.port = defaults.port;
    	      this.privateKey = defaults.privateKey;
    	      this.privateKeyPassword = defaults.privateKeyPassword;
    	      this.proxy = defaults.proxy;
    	      this.user = defaults.user;
        }

        @CustomType.Setter
        public Builder agentSocketPath(@Nullable String agentSocketPath) {
            this.agentSocketPath = agentSocketPath;
            return this;
        }
        @CustomType.Setter
        public Builder dialErrorLimit(@Nullable Integer dialErrorLimit) {
            this.dialErrorLimit = dialErrorLimit;
            return this;
        }
        @CustomType.Setter
        public Builder host(String host) {
            this.host = Objects.requireNonNull(host);
            return this;
        }
        @CustomType.Setter
        public Builder password(@Nullable String password) {
            this.password = password;
            return this;
        }
        @CustomType.Setter
        public Builder port(@Nullable Double port) {
            this.port = port;
            return this;
        }
        @CustomType.Setter
        public Builder privateKey(@Nullable String privateKey) {
            this.privateKey = privateKey;
            return this;
        }
        @CustomType.Setter
        public Builder privateKeyPassword(@Nullable String privateKeyPassword) {
            this.privateKeyPassword = privateKeyPassword;
            return this;
        }
        @CustomType.Setter
        public Builder proxy(@Nullable ProxyConnection proxy) {
            this.proxy = proxy;
            return this;
        }
        @CustomType.Setter
        public Builder user(@Nullable String user) {
            this.user = user;
            return this;
        }
        public Connection build() {
            final var o = new Connection();
            o.agentSocketPath = agentSocketPath;
            o.dialErrorLimit = dialErrorLimit;
            o.host = host;
            o.password = password;
            o.port = port;
            o.privateKey = privateKey;
            o.privateKeyPassword = privateKeyPassword;
            o.proxy = proxy;
            o.user = user;
            return o;
        }
    }
}
