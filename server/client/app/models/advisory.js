import DS from "ember-data";

export default DS.Model.extend({

  user: DS.belongsTo('user'),
  updates: DS.hasMany('update'),
  comments: DS.attr('raw'),
  description: DS.attr('string'),
  references: DS.attr('string'),
  status: DS.attr('number'),
  advisoryId: DS.attr('string'),
  createdAt: DS.attr('date'),
  updatedAt: DS.attr('date')

});
