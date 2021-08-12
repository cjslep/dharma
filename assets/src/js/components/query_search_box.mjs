export default {
	props: ['query'],
	emits: {
		updateQuery: function(value) {
			return value.length > 3 || value.length === 0;
		}
	},
	template: `<input :value="query" @input="$emit('updateQuery', $event.target.value)" />`
};
