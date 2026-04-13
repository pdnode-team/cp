<script lang="ts">
  import { page } from '$app/state'
  let modalRef: HTMLDialogElement
</script>

<div class="hero min-h-screen bg-base-100">
  <div class="hero-content text-center">
    <div class="max-w-lg">
      <h1 class="text-9xl font-bold">{page.status}</h1>
      <p class="py-6 text-xl">{page.error?.message ?? 'An unexpected error occurred.'}</p>
      <p class="py-6 text-lg">If you believe this issue should not exist, <a onclick={() => modalRef.showModal()} class="link link-primary" href="#contact-us">please contact us.</a></p>
      <a href="/" class="btn btn-primary">Back to Home</a>
    </div>
  </div>
</div>

<dialog bind:this={modalRef} class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Contact Us</h3>
    <p class="py-4 font-bold">You can contact us in a variety of ways.</p>

    <ul class="gap-2 flex flex-col">
        <li>Email: bret@pdnode.com</li>
        <li><a href="https://github.com/pdnode-team/cp/issues" class="link link-primary" target="_blank" rel="noopener noreferrer">GitHub Issues</a></li>
        <li><a href="https://chat.p67.click" class="link link-primary" target="_blank" rel="noopener noreferrer">Use Pdnode Team Chat (Zulip)</a></li>
    </ul>
    
    <div class="modal-action">
      <form method="dialog">
        <button class="btn">Close</button>
      </form>
    </div>
  </div>
</dialog>