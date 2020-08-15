<template>
  <div :id="'user-' + useruid">
    <b-card
      class="m-1 p-0 text-left"
      border-variant="dark"
      :header-bg-variant="stateHeaderBg"
      :header-text-variant="stateHeaderText"
    >
      <template v-slot:header>
        <b-row>
          <b-col class="text-left">
            <!-- UUID -->
            <label class="mb-2 mr-sm-2 mb-sm-0">UUID: {{ useruid }}</label>
          </b-col>
          <b-col class="text-right">
            <b-button size="sm" variant="dark" @click="collapse = !collapse">
              <b-icon-arrow-bar-up v-if="!collapse" />
              <b-icon-arrow-bar-down v-else />
            </b-button>
          </b-col>
        </b-row>
      </template>
      <b-form @submit="submit" @reset="reset" v-if="show && !collapse">
        <!-- Username -->
        <b-form-group id="input-group-username" label="Username:" label-for="input-username">
          <b-form-input
            id="input-username"
            v-model="form.username"
            type="text"
            required
            :state="stateUsername"
          />
        </b-form-group>

        <hr />
        <!-- Is admin -->
        <label class="mr-sm-2" for="inline-form-input-is-admin">Admin status</label>
        <b-form-checkbox
          class="mb-2 mr-sm-2 mb-sm-0"
          name="inline-form-input-is-admin"
          switch
          v-model="form.is_admin"
          :state="stateIsAdmin"
        />

        <hr />
        <!-- Buttons -->
        <b-button type="reset" variant="white" :disabled="!changedAny">
          <b-icon-x-square-fill />
        </b-button>

        <b-button type="submit" variant="white" :disabled="!changedAny">
          <b-icon-check2-square />
        </b-button>
      </b-form>
    </b-card>
  </div>
</template>

<script>
import errhandler from '../../utility/errhandler.js'

export default {
  name: "admin-user",

  data() {
    return {
      collapse: true,
      show: true,
      form: {
        username: this.username,
        is_admin: this.isAdmin,
      },
    };
  },

  computed: {
    stateHeaderText: function () {
      return this.changedAny ? "dark" : "white";
    },
    stateHeaderBg: function () {
      return this.changedAny ? "white" : "dark";
    },
    changedAny: function () {
      return this.changedUsername || this.changedIsAdmin;
    },

    stateUsername: function () {
      return this.changedUsername ? true : null;
    },
    changedUsername: function () {
      return this.username !== this.form.username;
    },

    stateIsAdmin: function () {
      return this.changedIsAdmin ? true : null;
    },
    changedIsAdmin: function () {
      return this.isAdmin !== this.form.is_admin;
    },
  },

  props: {
    useruid: { type: String, reqiuired: true },
    username: { type: String, reqiuired: true },
    created: { type: String, reqiuired: true },
    edited: { type: String, reqiuired: true },
    isAdmin: { type: Boolean, reqiuired: true },
  },

  methods: {
    reset(event) {
      event.preventDefault();
      this.form.username = this.username;
      this.form.is_admin = this.isAdmin;

      this.show = false;
      this.$nextTick(() => (this.show = true));
    },

    submit(event) {
      event.preventDefault();
      this.edit({ username: this.form.username, is_admin: this.form.is_admin });
    },

    edit(data) {
      this.$http({
        url: `admin/user/${this.useruid}`,
        data: data,
        method: "PATCH",
      })
        .then((response) => {
          this.form.username = response.data.username;
          this.$emit("update:username", response.data.username);

          this.form.is_admin = response.data.is_admin;
          this.$emit("update:isAdmin", response.data.is_admin);

          this.created = response.data.created;
          this.edited = response.data.edited;
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = 'Admin users editing error' + (code ? ` with code ${code}` : "");
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

<style>
</style>