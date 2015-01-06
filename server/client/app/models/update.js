import DS from "ember-data";

export default DS.Model.extend({

  advisory: DS.belongsTo('advisory'),
  for: DS.attr('string'),
  name: DS.attr('string'),
  submitter: DS.attr('string'),
  type: DS.attr('string'),
  available: DS.attr('boolean'),
  createdAt: DS.attr('date'),
  content: DS.attr('raw')

});
