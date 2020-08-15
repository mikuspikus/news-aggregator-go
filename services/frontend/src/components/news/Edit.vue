<template>
  <div id="news-edit-form">
    <b-form @submit="submit" @reset="reset" v-if="show">
      <b-form-group id="input-group-news-title" label="Title:" label-for="input-news-title">
        <b-form-input
          id="input-news-title"
          v-model="form.title"
          type="text"
          placeholder="Enter title"
          required
        />
      </b-form-group>

      <b-form-group id="input-group-new-url" label="URL:" label-for="input-news-url">
        <b-form-input
          id="input-news-url"
          v-model="form.uri"
          type="url"
          placeholder="Enter URL"
          required
        />
      </b-form-group>

      <div id="form-errors">
        <b-alert show dismissible variant="light" v-for="error in errors" :key="error.message">
          <hr />
          <p>{{ error.message}}</p>
          <hr />
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
import errhandler from '../../utility/errhandler.js'

export default {
  name: "news-edit-form",

  props: {
    uid: { type: String, required: true },
    title: { type: String, required: true },
    uri: { type: String, required: true },
    edited: { type: String, required: true },
    edit: { type: Boolean, required: true },
  },

  data() {
    return {
      show: true,
      errors: [],
      form: {
        title: this.title,
        uri: this.uri,
      },
    };
  },

  methods: {
    update(uid, data) {
      this.$http({
        url: `news/${uid}`,
        data: data,
        method: "PATCH",
      })
        .then((response) => {
          this.$emit("update:title", response.data.title);
          this.$emit("update:uri", response.data.uri);
          this.$emit("update:edited", response.data.edited);
          this.$emit("update:edit", false);
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "News editing error" + (code ? ` with code ${code}` : "");
          
          this.$bvToast.toast(message, {
            title: title,
            autoHideDelay: 5000,
            variant: "white",
            toaster: "b-toaster-bottom-center",
          });
        });
    },

    submit(event) {
      event.preventDefault();

      this.update(this.uid, this.form);
    },

    reset(event) {
      event.preventDefault();

      this.form.title = this.title;
      this.form.uri = this.uri;
      this.show = false;
      this.$nextTick(() => (this.show = true));
    },
  },
};
</script>