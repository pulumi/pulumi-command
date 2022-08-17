// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.command.remote.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import java.lang.Double;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Instructions for how to connect to a remote endpoint.
 * 
 */
public final class ConnectionArgs extends com.pulumi.resources.ResourceArgs {

    public static final ConnectionArgs Empty = new ConnectionArgs();

    /**
     * The address of the resource to connect to.
     * 
     */
    @Import(name="host", required=true)
    private Output<String> host;

    /**
     * @return The address of the resource to connect to.
     * 
     */
    public Output<String> host() {
        return this.host;
    }

    /**
     * The password we should use for the connection.
     * 
     */
    @Import(name="password")
    private @Nullable Output<String> password;

    /**
     * @return The password we should use for the connection.
     * 
     */
    public Optional<Output<String>> password() {
        return Optional.ofNullable(this.password);
    }

    /**
     * The port to connect to.
     * 
     */
    @Import(name="port")
    private @Nullable Output<Double> port;

    /**
     * @return The port to connect to.
     * 
     */
    public Optional<Output<Double>> port() {
        return Optional.ofNullable(this.port);
    }

    /**
     * The contents of an SSH key to use for the connection. This takes preference over the password if provided.
     * 
     */
    @Import(name="privateKey")
    private @Nullable Output<String> privateKey;

    /**
     * @return The contents of an SSH key to use for the connection. This takes preference over the password if provided.
     * 
     */
    public Optional<Output<String>> privateKey() {
        return Optional.ofNullable(this.privateKey);
    }

    /**
     * The user that we should use for the connection.
     * 
     */
    @Import(name="user")
    private @Nullable Output<String> user;

    /**
     * @return The user that we should use for the connection.
     * 
     */
    public Optional<Output<String>> user() {
        return Optional.ofNullable(this.user);
    }

    private ConnectionArgs() {}

    private ConnectionArgs(ConnectionArgs $) {
        this.host = $.host;
        this.password = $.password;
        this.port = $.port;
        this.privateKey = $.privateKey;
        this.user = $.user;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ConnectionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ConnectionArgs $;

        public Builder() {
            $ = new ConnectionArgs();
        }

        public Builder(ConnectionArgs defaults) {
            $ = new ConnectionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param host The address of the resource to connect to.
         * 
         * @return builder
         * 
         */
        public Builder host(Output<String> host) {
            $.host = host;
            return this;
        }

        /**
         * @param host The address of the resource to connect to.
         * 
         * @return builder
         * 
         */
        public Builder host(String host) {
            return host(Output.of(host));
        }

        /**
         * @param password The password we should use for the connection.
         * 
         * @return builder
         * 
         */
        public Builder password(@Nullable Output<String> password) {
            $.password = password;
            return this;
        }

        /**
         * @param password The password we should use for the connection.
         * 
         * @return builder
         * 
         */
        public Builder password(String password) {
            return password(Output.of(password));
        }

        /**
         * @param port The port to connect to.
         * 
         * @return builder
         * 
         */
        public Builder port(@Nullable Output<Double> port) {
            $.port = port;
            return this;
        }

        /**
         * @param port The port to connect to.
         * 
         * @return builder
         * 
         */
        public Builder port(Double port) {
            return port(Output.of(port));
        }

        /**
         * @param privateKey The contents of an SSH key to use for the connection. This takes preference over the password if provided.
         * 
         * @return builder
         * 
         */
        public Builder privateKey(@Nullable Output<String> privateKey) {
            $.privateKey = privateKey;
            return this;
        }

        /**
         * @param privateKey The contents of an SSH key to use for the connection. This takes preference over the password if provided.
         * 
         * @return builder
         * 
         */
        public Builder privateKey(String privateKey) {
            return privateKey(Output.of(privateKey));
        }

        /**
         * @param user The user that we should use for the connection.
         * 
         * @return builder
         * 
         */
        public Builder user(@Nullable Output<String> user) {
            $.user = user;
            return this;
        }

        /**
         * @param user The user that we should use for the connection.
         * 
         * @return builder
         * 
         */
        public Builder user(String user) {
            return user(Output.of(user));
        }

        public ConnectionArgs build() {
            $.host = Objects.requireNonNull($.host, "expected parameter 'host' to be non-null");
            $.port = Codegen.doubleProp("port").output().arg($.port).def(2.2e+01).getNullable();
            $.user = Codegen.stringProp("user").output().arg($.user).def("root").getNullable();
            return $;
        }
    }

}
