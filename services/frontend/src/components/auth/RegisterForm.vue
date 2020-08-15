<template>
  <div id="register-form">
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

      <b-form-group id="input-group-password1" label="Password:" label-for="input-password1">
        <b-form-input
          id="input-password1"
          v-model="form.password1"
          type="password"
          placeholder="Enter your password"
          required
        />
      </b-form-group>

      <b-form-group
        id="input-group-password2"
        label="Password (second time):"
        label-for="input-password2"
      >
        <b-form-input
          id="input-password2"
          v-model="form.password2"
          type="password"
          placeholder="Enter your password (second time)"
          required
        />
      </b-form-group>

      <div id="form-errors">
        <b-alert show dismissible variant="light" v-for="error in errors" :key="error.message">
          <h4 class="alert-heading">Error!</h4>
          <hr />
          <p>{{ error.message}}</p>
        </b-alert>
      </div>

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
  name: "register-form",

  data() {
    return {
      show: true,
      errors: [],
      form: {
        username: "",
        password1: "",
        password2: "",
      },
    };
  },

  methods: {
    reset(event) {
      event.preventDefault();

      this.form.username = "";
      this.form.password1 = "";
      this.form.password2 = "";
      // this.errors = [];

      this.show = false;
      this.$nextTick(() => (this.show = true));
    },

    submit(event) {
      event.preventDefault();

      this.errors = [];

      if (this.form.password1 !== this.form.password2) {
        this.errors.push({ message: "Your passwords do not match" });
        return;
      }

      this.$store
        .dispatch("register", {
          credentials: {
            username: this.form.username,
            password: this.form.password1,
          },
          axios: this.$http,
        })
        .then(() => this.$router.push({ name: "Login" }))
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "Register error" + (code ? ` with code ${code}` : "");

          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },
  },
};
</script>