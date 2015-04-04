import Ember from "ember";

export default Ember.Route.extend({
  titleToken: "Update",
  model: function(params) {
    return Ember.$.getJSON('/api/v1/updates/' + params.id);
  },
  afterModel: function(posts, transition) {
    this.titleToken = this.model['update']['name'];
  }
});
