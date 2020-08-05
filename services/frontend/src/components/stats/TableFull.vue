<template>
  <div id="stats-full">
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
            v-model="currentPage"
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
      <b-table :items="items" :fields="fields" stacked="md" striped small bordered dark>
        <template v-slot:cell(io)="row">
          <b-row>
            <b-col>
              <b-row>
                <b-col>Input:</b-col>
                <b-col>{{ row.item.input.size }}</b-col>
              </b-row>

              <b-row>
                <b-col>Output:</b-col>
                <b-col>{{ row.item.output.size }}</b-col>
              </b-row>
            </b-col>
            <b-col align-self="center" md="auto" class="m-2">
              <b-button
                v-if=" row.item.input.size || row.item.output.size"
                size="sm"
                variant="dark"
                @click="row.toggleDetails"
              >
                <b-icon-arrow-bar-up v-if="row.detailsShowing" />
                <b-icon-arrow-bar-down v-else />
              </b-button>
            </b-col>
          </b-row>
        </template>

        <template v-slot:row-details="row">
          <b-card class="text-left" bg-variant="dark">
            <b-row v-if="row.item.input.size">
              <json-view
                :maxDepth="1"
                :rootKey="'Input'"
                :data="row.item.input.data"
                :colorScheme="'dark'"
              />
            </b-row>

            <b-row v-if="row.item.output.size">
              <json-view
                :maxDepth="1"
                :rootKey="'Output'"
                :data="row.item.output.data"
                :colorScheme="'dark'"
              />
            </b-row>
          </b-card>
        </template>
      </b-table>
    </b-container>
  </div>
</template>

<script>
import { JSONView } from "vue-json-component";
export default {
  name: "stats-panel-full",

  components: { "json-view": JSONView },

  created() {
    this.fetch(this.currentPage, this.pageSize);
  },

  methods: {
    _processStats(item) {
      const tmp = item;
      tmp.input = {
        data: tmp.input,
        size: this._objSize(tmp.input),
      };
      tmp.output = {
        data: tmp.output,
        size: this._objSize(tmp.output),
      };
      this.items.push(tmp);
    },

    changePage(newPage) {
      this.currentPage = newPage;
      this.fetch(this.currentPage, this.pageSize);
    },

    changePageSize(newPageSize) {
      this.pageSize = newPageSize;
      this.changePage(1);
    },

    _objSize(obj) {
      return obj ? Object.keys(obj).length : 0;
    },

    fetch(page, size) {
      this.items = [];
      console.log(page, size);
      this.$http({
        url: "stats/news",
        params: { page: page - 1, size: size },
        method: "GET",
      })
        .then((response) => {
          console.log(response.data);
          response.data.stats.forEach(this._processStats);
          this.pageCount = response.data.page_count;
        })
        .catch();
    },
  },

  data() {
    return {
      items: [],
      fields: [
        {
          key: "id",
          label: "ID",
          sortable: true,
          sortDirection: "desc",
          class: "align-middle",
        },
        {
          key: "user",
          label: "User",
          sortable: true,
          class: "align-middle",
          formatter: (value) => {
            return value ? value : "Anonymous";
          },
        },
        {
          key: "timestamp",
          label: "Timestamp",
          sortable: true,
          formatter: (value /*, key, item*/) => {
            return this.$moment(value).format("dddd, MMMM Do YYYY, HH:mm:ss");
          },
          class: "align-middle",
        },
        {
          key: "action",
          label: "Action",
          sortable: true,
          class: "align-middle",
        },
        { key: "io", label: "Input/Output", class: "align-middle" },
      ],
      currentPage: 1,
      pageSize: 5,
      pageCount: 1,
      pageOptions: [5, 10, 15, 20, 25],
    };
  },
};
</script>

<style>
.pagination > li > a {
  background-color: white;
  color: #343a40;
}

.pagination > li > a:focus,
.pagination > li > a:hover,
.pagination > li > span:focus,
.pagination > li > span:hover {
  color: #343a40;
  background-color: #eee;
  border-color: #ddd;
}

.pagination > .active > a {
  color: white;
  background-color: #343a40 !important;
  border: solid 1px #343a40 !important;
}

.pagination > .active > a:hover {
  background-color: #343a40 !important;
  border: solid 1px #343a40;
}
</style>