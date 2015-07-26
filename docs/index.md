<div class="jumbotron intro">
Dcron is a system service that runs scheduled tasks at given intervals or times, just like the cron unix service. It differs from it in the sense that it's distributed in several machines in a cluster and if one of that machines (the leader) fails, any other one can take this responsability and keep executing the sheduled tasks without human intervention.
</div>

## Characteristics

<div id="easy-integration" class="row vertical-align">
  <img src="img/integration.png" class="col-md-3"/>
  <div class="col-md-10">
    <h3>Easy integration</h3>
    Dcron is easy to setup and use. Choose your OS package and it's ready to run out-of-the-box. The [administration panel]() and it's simple JSON API makes a breeze to integrate with you current workflow or deploy system.
  </div>
</div>

<div id="always-available" class="row vertical-align">
  <div class="col-md-10">
    <h3>Always available</h3>
    Using the power of the Raft implementation in etcd, you can rely on Dcron to be always available. Wherever you need to generate your employee monthly payroll or send those daily emails that keeps your users informed and your site busy, Dcron could help you to sleep at night.
  </div>
  <img src="img/available.png" class="col-md-3"/>
</div>

<div id="flexible-targets" class="row vertical-align">
  <img src="img/targets.png" class="col-md-3"/>
  <div class="col-md-10">
    <h3>Flexible targets</h3>
    Simple but powerful tag based target node selection for jobs. Tag node count allows to run jobs in an arbitrary number nodes in the same group or groups for any given task.
  </div>
</div>
## Example use cases:

* Eail delivery: Newsletters, marketing campaigns, recommendations, etc.
* Payroll generation
* Bookkeeping
* Data consolidation for BI
* Recurring invoicing
* Data transfer
* ...
