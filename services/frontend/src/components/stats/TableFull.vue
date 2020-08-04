<template>
  <div id="stats-full">
    <b-table
      :items="items"
      :fields="fields"
      stacked="md"
      :current-page="currentPage"
      striped
      small
      bordered
      dark
    >
      <template v-slot:cell(io)="row">
        <b-row>
          <b-col>
            <b-row>
              <b-col>Input:</b-col>
              <b-col>{{ _size(row.item.input) }}</b-col>
            </b-row>

            <b-row>
              <b-col>Output:</b-col>
              <b-col>{{ _size(row.item.output) }}</b-col>
            </b-row>
          </b-col>
          <b-col align-self="center" md="auto" class="m-2">
            <b-button
              v-if=" _size(row.item.input) || _size(row.item.output)"
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
          <b-row v-if="_size(row.item.input)">
            <json-view
              :maxDepth="1"
              :rootKey="'Input'"
              :data="row.item.input"
              :colorScheme="'dark'"
            />
          </b-row>

          <b-row v-if="_size(row.item.output)">
            <json-view
              :maxDepth="1"
              :rootKey="'Output'"
              :data="row.item.output"
              :colorScheme="'dark'"
            />
          </b-row>
        </b-card>
      </template>
    </b-table>
  </div>
</template>

<script>
import { JSONView } from "vue-json-component";
export default {
  name: "stats-panel-full",

  components: { "json-view": JSONView },

  methods: {
    _size(obj) {
      return Object.keys(obj).length;
    },
  },

  data() {
    return {
      items: [
        {
          id: 1,
          user: "",
          action: "List News",
          timestamp: "2013-02-08 24:00:00.000",
          input: {
            "wubba-laba-dub-dub": "pickle rick",
            "asdfhklashfasf": {'aslkdnfhlsakhdflkasf': 'dsafphslfdnhkhas'},
          },
          output: {},
        },
      ],
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
      pageOptions: [5, 10, 15, 20, 25],
    };
  },
};
</script>

<style>
</style>