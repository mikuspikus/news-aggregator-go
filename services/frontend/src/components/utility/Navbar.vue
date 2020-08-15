<template>
  <div id="main-navbar">
    <b-navbar toggleable="lg" type="dark" variant="dark">
      <b-navbar-brand :to="{ name: 'Home'}">News Aggregator Go</b-navbar-brand>
      <b-navbar-toggle target="nav-collapse" />

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item :to="{ name: 'Home' }">Home</b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav class="ml-auto">
          <b-nav-item v-if="isAdmin" :to="{ name: 'AdminPanel'}">Admin panel</b-nav-item>

          <b-nav-item v-if="isLogged" @click="logout">Sign out</b-nav-item>

          <template v-else>
            <b-nav-item :to="{ name: 'Login' }">Sign in</b-nav-item>
            <b-nav-item :to="{ name: 'Register' }">Sign up</b-nav-item>
          </template>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
  </div>
</template>

<script>
import errhandler from "../../utility/errhandler.js";

export default {
  name: "Navbar",

  props: {
    isLogged: { type: Boolean, required: true },
    isAdmin: { type: Boolean, required: true },
  },

  methods: {
    logout(event) {
      event.preventDefault();

      this.$store
        .dispatch("logout", { axios: this.$http })
        .then(() => {
          this.$router.push({ name: "Home" }).catch(() => {});
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "Logout error" + (code ? ` with code ${code}` : "");
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-right",
          });
        });
    },
  },
};
</script>