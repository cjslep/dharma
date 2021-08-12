import SearchResultList from '../components/search_result_list.mjs';

export default {
	data: function() {
		return {
			query: '',
			resultList: [],
		};
	},
	props: {
		searchText: String,
		submitText: String,
		searchEndpoint: String
	},
	components: {
		'search-result-list': SearchResultList
	},
	computed: {
		isValidQuery: function() {
			console.log(this.query)
			return this.query.length > 3
		},
	},
	methods: {
		onSearch: function(event) {
			let data = new FormData();
			data.append('query', this.query);
			const opts = {
				method: "POST",
				body: data
			};
			fetch(this.searchEndpoint, opts)
				.then(response => response.json())
				.then(data => this.onSuccessJSON(data))
				.catch(error => this.onSearchError(error));
		},
		onSuccessJSON: function(data) {
			this.resultList = data.corporations;
		},
		onSearchError: function(error) {
			console.error(error);
		}
	},
	template: `
<div>
  <input v-model="query" type="text" />
  <button :disabled="!isValidQuery" @click="onSearch">{{searchText}}</button>
  <search-result-list :result-list="this.resultList" :submit-text="submitText" />
</div>`
}
