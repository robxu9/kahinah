import Ember from "ember";

export default Ember.Route.extend({
  titleToken: "Latest Updates",
  model: function() {
    return Ember.RSVP.hash({
      updates: Ember.$.getJSON("/api/v1/updates"),
      targets: Ember.$.getJSON("/api/v1/updates/targets")
    });
  }
});
