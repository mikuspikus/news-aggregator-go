<template>
    <div id="news-short">
        <b-card class="m-0 p-0 text-left">
            <b-card-text class="m-0 p-0">
                <b-link :to="{ name: 'SingleNews', params: { uid: uid } }">{{title}}</b-link>
            </b-card-text>
            <template v-slot:footer>
                <small class="text-muted p-0 m-0">Created at {{created | moment("dddd, MMMM Do YYYY, HH:mm:ss")}}</small>
            </template>
        </b-card>
    </div>
</template>

<script>
export default {
    name: "news-short",

    props: {
        uid: { type: String, required: true },
        title: { type: String, required: true },
        created: { type: String, required: true },
    }
}
</script>