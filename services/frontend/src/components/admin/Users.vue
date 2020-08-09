<template>
  <div id="admin-users">
    <b-container fluid>
      <b-row>
        <b-col>
          <b-form-group
            label="Per page"
            label-cols-sm="6"
            label-cols-md="4"
            label-cols-lg="3"
            label-align-sm="right"
            label-size="sm"
            label-for="perPageSelect"
            class="mb-0"
          >
            <b-form-select
              v-model="pageSize"
              id="perPageSelect"
              size="sm"
              :options="pageOptions"
              @input="changePageSize"
            />
          </b-form-group>
        </b-col>

        <b-col sm="7" md="6">
          <b-pagination
            v-model="currnetPageNumber"
            @input="changePage"
            :total-rows="pageCount"
            :per-page="1"
            first-number
            last-number
            align="fill"
            size="sm"
          />
        </b-col>
      </b-row>

      <!-- <b-row > -->
      <user
        v-for="user in users"
        :key="user.uid"
        :useruid="user.uid"
        :username.sync="user.username"
        :created="user.created"
        :edited="user.edited"
        :isAdmin.sync="user.is_admin"
      />
      <!-- </b-row> -->
    </b-container>
  </div>
</template>

<script>
import User from "../admin/User.vue";

export default {
  name: "admin-users",

  components: { User },

  data() {
    return {
      users: [],

      currnetPageNumber: 1,
      pageSize: 25,
      pageOptions: [5, 10, 15, 20, 25, 50],
      pageCount: 1,
    };
  },

  created() {
    this.fetch(this.currnetPageNumber, this.pageSize);
  },

  methods: {
    changePage(newPageNumber) {
      this.currnetPageNumber = newPageNumber;
      this.fetch(this.currnetPageNumber, this.pageSize);
    },

    changePageSize(newPageSize) {
      this.pageSize = newPageSize;
      this.changePage(1);
    },

    fetch(pageNumber, pageSize) {
      this.users = [];
      this.$http({
        url: "admin/user/",
        params: { page: pageNumber - 1, size: pageSize },
        method: "GET",
      })
        .then((response) => {
          this.users = response.data.users;
          this.pageCount = response.data.page_count;
        })
        .catch((error) => {
          this.$bvToast.toast(error, {
            title: "Admin users fetching error",
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