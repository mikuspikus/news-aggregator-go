<template>
    <div id="service-error-card">
        <b-jumbotron
            bg-variant="dark"
            text-variant="white"
            border-variant="dark"
            class="m-0"
            header
        >
            <template v-slot:header>{{ header }}</template>
            <template v-slot:lead>{{ lead }}</template>
        </b-jumbotron>
    </div>
</template>

<script>
export default {
    name: "service-error-card",

    props: {
        header: { type: String, default: "Error occuried" },
        lead: { type: String, default: "Some service is probably unavailible" }
    }
}
</script>