<template>
  <div id="register">
    <b-container fluid>
      <b-row class="text-left">
        <b-col cols="4"></b-col>
        <b-col class="pt-5 pb-5">
          <b-card
            border-variant="dark"
            header="Sign up"
            header-text-variant="white"
            header-bg-variant="dark"
            header-class="text-center"
          >
            <register-form />
          </b-card>
        </b-col>
        <b-col cols="4"></b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import RegisterForm from "@/components/auth/RegisterForm.vue";

export default {
  name: "register-view",

  components: {
    RegisterForm,
  },
};
</script>