import {app} from '../base/vue.mjs';
import CorporationSearch from '../components/corporation_search.mjs';

app.component('corporation-search', CorporationSearch)
	.mount('#dharma-app');
