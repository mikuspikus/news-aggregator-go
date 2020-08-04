<template>
  <div id="admin-panel">
    <b-container fluid class="mt-2 mb-2">
      <b-card no-body>
        <b-tabs content-class="mt-3" no-fade lazy>
          <b-tab title="First" active>
            <stats-table-full />
          </b-tab>
          <b-tab title="Second">
            <p>I'm the second tab</p>
          </b-tab>
          <b-tab title="Disabled">
            <p>I'm a disabled tab!</p>
          </b-tab>
        </b-tabs>
      </b-card>
    </b-container>
  </div>
</template>

<script>
import StatsTableFull from "../components/stats/TableFull.vue";
export default {
  name: "admin-panel",

  components: { StatsTableFull },
};
</script>