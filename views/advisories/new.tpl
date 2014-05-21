{{template "header.tpl" .}}

		<div class="page-header">
			<h1>Create Advisory <small>Something big needs to be known.</small></h1>
		</div>

		<form class="form-horizontal" role="form" id="advisory_new" method="post">
			<div class="form-group"> <!-- SUMMARY -->
				<label for="input_summary" class="col-sm-2 control-label">Summary</label>
				<div class="col-sm-10">
			    	<input type="text" class="form-control" id="input_summary" placeholder="What's this about?"{{if .FailSummary}}value="{{.FailSummary}}"{{end}} name="input_summary">
			    </div>
			</div>
			<div class="form-group"> <!-- PLATFORM -->
				<label for="input_platform" class="col-sm-2 control-label">Platform</label>
				<div class="col-sm-10">
					<select id="input_platform" class="form-control" name="input_platform">
						{{range $key, $value := .Platforms}}
						<option value="{{$key}}">{{$key}}</option>
						{{end}}
					</select>
  				</div>
			</div>
			<div class="form-group"> <!-- DESCRIPTION -->
				<label for="input_description" class="col-sm-2 control-label">Description</label>
				<div class="col-sm-10">
					<textarea class="form-control" rows="3" id="input_description" placeholder="Explain, in detail, what this update encompasses. Include CVE numbers and other details." form="advisory_new" name="input_description">{{if .FailDescription}}{{.FailDescription}}{{end}}</textarea>
			    </div>
			</div>
			<div class="form-group"> <!-- TYPE -->
				<label for="input_type" class="col-sm-2 control-label">Type</label>
				<div class="col-sm-10">
					<select id="input_type" class="form-control" name="input_type">
						{{range $key, $value := .Types}}
						<option>{{$value}}</option>
						{{end}}
					</select>
			    </div>
			</div>
			<div class="form-group"> <!-- BUGS FIXED -->
				<label for="input_bugs" class="col-sm-2 control-label">Bugs Fixed <button class="btn btn-link" id="bugs_add">Add</button></label>
				<div class="col-sm-10" id="bugs_div">
					{{if .FailBug}}{{range $key, $value := .FailBug}}<input type="number" autocomplete="off" class="form-control input_bugs" name="input_bugs" id="input_bugs_{{$key}}" placeholder="Update ID" value="{{$value}}" min="0"/>{{end}}{{else}}<input type="number" autocomplete="off" class="form-control input_bugs" name="input_bugs" id="input_bugs_0" placeholder="e.g. 789 (for omv#789)" min="0"/>{{end}}
				</div>
			</div>
			<div class="form-group"> <!-- UPDATE IDs -->
				<label for="input_update" class="col-sm-2 control-label">Update IDs <button class="btn btn-link" id="update_add">Add</button></label>
				<div class="col-sm-10" id="update_div">
					{{if .FailUpdateID}}{{range $key, $value := .FailUpdateID}}<input type="number" autocomplete="off" class="form-control input_update" name="input_update" id="input_update_{{$key}}" placeholder="Update ID" value="{{$value}}" min="0"/>{{end}}{{else}}<input type="number" autocomplete="off" class="form-control input_update" name="input_update" id="input_update_0" placeholder="e.g. 152 (for UPDATE-YEAR-152)" min="0"/>{{end}}
				</div>
			</div>
			<div class="form-group"> <!-- SUBMIT -->
				<label for="input_submit" class="col-sm-6 control-label">Make sure to proofread! Once submitted, this cannot be changed!</label>
				<div class="col-sm-6">
					<input type="submit" class="btn btn-primary" value="I have proofread this advisory, and know that this advisory can not be edited.">
				</div>
			</div>
			{{.xsrf_data}}
		</form>

		<script>
			$(document).ready(function() { 
				var update_i = $(".input_update").length;

				$("#update_add").click(function(e) {
					e.preventDefault();
				    $("#input_update_0").clone()
				        .attr("id", "input_update_" + update_i)
				        .val("")
				        .appendTo("#update_div");
				    update_i++;
				});

				var bugs_i = $(".input_bugs").length;

				$("#bugs_add").click(function(e) {
					e.preventDefault();
				    $("#input_bugs_0").clone()
				        .attr("id", "input_bugs_" + bugs_i)
				        .val("")
				        .appendTo("#bugs_div");
				    bugs_i++;
				});
			});
		</script>

{{template "footer.tpl" .}}