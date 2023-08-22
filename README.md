# circlepipe
Generate circleci pipeline for use with the continuation orb



Customizations and Defaults


list-of-the-roles-pre-jobs = ""
approval-jobs = ""

For role in roles


  # pre section
  if prerole-only
    add job to the list of pre-jobs
    fetch compiledEnv
    generate pre template and add to .yml
    requires approval-job  (which is blank the first time through)

  else
    for instance in instances
      add job to list of pre-jobs
      fetch compiledEnv
      generate pre template and add to .yml

  Generate "approval" job
    set approval-job name based on role
    requires -list-of-the-roles-pre-jobs

  clear list-of-the-roles-pre-jobs

  # post section
  if postrole-only
    add job to the list of pre-jobs
    fetch compiledEnv
    generate post template and add to .yml
    requires approval-job

  else
    for instance in instances
      add job to list of pre-jobs
      fetch compiledEnv
      generate pre template and add to .yml
