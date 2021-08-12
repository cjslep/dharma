export default {
	props: {
		resultList: Array,
		submitText: String,
	},
	template: `<div>
<form v-for="(value) in resultList" method="post">
	<p>{{value.name}}, [{{value.ticker}}]: {{value.ceo.name}}</p>
	<input type="hidden" name="corporation_id" :value="value.id" />
	<button>{{submitText}}</button>
</form>
</div>`
}
