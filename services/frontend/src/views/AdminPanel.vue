<template>
  <div id="admin-panel">
    <b-container fluid class="mt-2 mb-2">
      <b-card no-body>
        <b-tabs content-class="mt-3" no-fade lazy>
          <b-tab title="'News' stats" active>
            <stats-table-full :serviceName="'news'" />
          </b-tab>
          <b-tab title="'Comments' stats">
            <stats-table-full :serviceName="'comments'" />
          </b-tab>
          <b-tab title="'Accounts' stats'">
            <stats-table-full :serviceName="'accounts'" />
          </b-tab>
          <b-tab title="Users management">
            <users />
          </b-tab>
        </b-tabs>
      </b-card>
    </b-container>
  </div>
</template>

<script>
import StatsTableFull from "../components/admin/TableFull.vue";
import Users from "../components/admin/Users.vue";
export default {
  name: "admin-panel-view",

  components: { StatsTableFull, Users },
};
</script>