<template>
  <div id="login">
    <b-container fluid>
      <b-row class="text-left">
        <b-col cols="4"></b-col>
        <b-col class="pt-5 pb-5">
          <b-card
            border-variant="dark"
            header-text-variant="white"
            header-bg-variant="dark"
            header-class="text-center"
            header="Sign in"
          >
            <login-form />
          </b-card>
        </b-col>
        <b-col cols="4"></b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import LoginForm from "@/components/auth/LoginForm.vue";

export default {
  name: "login-view",

  components: {
    LoginForm,
  },
};
</script>