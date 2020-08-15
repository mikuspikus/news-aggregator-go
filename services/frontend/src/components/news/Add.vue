<template>
  <div id="add-news-form">
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

      <b-form-group id="input-group-new-url" label="URL:" label-for="input-newes-url">
        <b-form-input
          id="input-news-url"
          v-model="form.url"
          type="url"
          placeholder="Enter URL"
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
import errhandler from '../../utility/errhandler.js'

export default {
  name: "add-news-form",

  data() {
    return {
      errors: [],
      show: true,
      form: {
        title: "",
        url: "",
      },
    };
  },

  methods: {
    post(data) {
      this.$http({
        url: "news/",
        data: data,
        method: "POST",
      })
        .then((response) => {
          this.$router.push({
            name: 'SingleNews',
            params: { uid: response.data.uid }
          })
        })
        .catch((error) => {
          const { message, code } = errhandler.handle(error);
          const title = "News adding error" + (code ? ` with code ${code}` : "");
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

      this.errors = []

      const news = {
        title: this.form.title,
        uri: this.form.url
      };

      this.post(news)
    },

    reset(event) {
      event.preventDefault();

      this.form.title = "";
      this.form.url = "";

      this.show = false;
      this.$nextTick(() => this.show = true )
    }
  },
};
</script>