
      <!-- Jumbotron -->
      <div class="jumbotron">
        <h1>Kahinah</h1>
        <p class="lead">... is a versatile tool that hooks into build systems, allowing package developers to focus on developing more without worrying about breakage.</p>
        <p><a class="btn btn-lg btn-success" href="{{url "/builds/testing"}}" role="button">Recent Builds</a> <a class="btn btn-lg btn-default" href="{{url "/advisories"}}" role="button">Advisories</a> <a class="btn btn-lg btn-warning" href="{{url "/vtests"}}" role="button">Virtual Testing</a></p>
      </div>

      <!-- Infos -->
      <div class="row">
        <div class="col-lg-6">
          <h2>Breakage is expected.</h2>
          <p class="text-danger">Even though there has been great usage of this tool, things are ever evolving. Take caution.</p>
          <p><a class="btn btn-primary" href="{{url "/about"}}" role="button">Contact &raquo;</a></p>
        </div>
        <div class="col-lg-6">
          <h2>News</h2>
          {{.News}}
       </div>
      </div>
