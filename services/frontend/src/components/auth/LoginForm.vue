<template>
  <div id="login-form">
    <b-form @submit="submit" @reset="reset" v-if="show">
      <b-form-group id="input-group-username" label="Username:" label-for="input-username">
        <b-form-input
          id="input-username"
          v-model="form.username"
          type="text"
          placeholder="Enter your username"
          required
        />
      </b-form-group>

      <b-form-group id="input-group-password" label="Password:" label-for="input-password">
        <b-form-input
          id="input-password"
          v-model="form.password"
          type="password"
          placeholder="Enter your password"
          required
        />
      </b-form-group>

      <b-button class="mr-1" type="reset" variant="white">
        <b-icon-x-square-fill />
      </b-button>

      <b-button class="ml-1" type="submit" variant="white">
        <b-icon-check2-square />
      </b-button>
    </b-form>
  </div>
</template>

<script>
import errhandler from "../../utility/errhandler.js";

export default {
  name: "login-form",

  data() {
    return {
      show: true,
      form: {
        username: "",
        password: "",
      },
    };
  },

  methods: {
    submit(event) {
      event.preventDefault();

      this.$store
        .dispatch("login", {
          credentials: {
            username: this.form.username,
            password: this.form.password,
          },
          axios: this.$http,
        })
        .then(() => this.$router.push({ name: "Home" }))
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "Login error" + (code ? ` with code ${code}` : "");
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    reset(event) {
      event.preventDefault();

      this.form.username = "";
      this.form.password = "";

      this.show = false;
      this.$nextTick(() => (this.show = true));
    },
  },
};
</script>