import DS from "ember-data";

export default DS.Model.extend({

  email: DS.attr('string'),
  advisories: DS.hasMany('advisory'),
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date')

});
