<script lang="ts">
	import pb from '$lib/pocketbase';
	import TagInput from './ui/TagInput.svelte';
	import { toFormData } from '$lib/utils/api';

	let { open = $bindable(false), afterCreate } = $props();
	let dialogRef = $state() as HTMLDialogElement;

	$effect(() => {
		if (open) {
			dialogRef?.showModal();
		} else {
			dialogRef?.close();
		}
	});

	let tags = $state<string[]>([]);
	let errorText = $state('');

	let name = $state('');
	let description = $state('');
	let origin = $state('');
	let pictures: FileList | undefined = $state();

	const createCharacter = async () => {
		errorText = '';

		if (!name.trim() || !description.trim() || !origin.trim()) {
			errorText = 'All fields must be filled out';
			return;
		}

		try {
			await pb
				.collection('characters')
				.create(toFormData({ name, description, origin, tag_names: tags, images: pictures, owner: pb.authStore.record!.id }));
			name = '';
			description = '';
			origin = '';
			pictures = undefined;
			tags = [];
			errorText = '';
			open = false;
			if (afterCreate) {
				afterCreate();
			}
		} catch (err: any) {
			errorText = err.data.data?.message ?? 'Create failed. Please try again.';

			const firstKey = Object.keys(err.data.data)[0];
			if (firstKey) {
				const friendlyMessage = `${firstKey}: ${err.data.data[firstKey].message}`;
				console.error(friendlyMessage);
				errorText = friendlyMessage;
				return;
			}
		}
	};
</script>

<dialog bind:this={dialogRef} onclose={() => (open = false)} class="modal">
	<div class="modal-box flex flex-col gap-4">
		<h3 class="text-lg font-bold">Create a new character</h3>
		<div class="flex flex-col gap-4">
			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Your Character Name</span></div>
				<input
					type="text"
					name="characterName"
					placeholder="XXXX & YYYY"
					class="input-bordered input w-full"
					autocomplete="name"
					bind:value={name}
					required
				/>
			</label>

			<label class="form-control w-full">
				<div class="label">
					<span class="label-text font-medium">Your Character Description</span>
				</div>
				<textarea
					name="characterDescription"
					placeholder="description"
					class="textarea w-full"
					maxlength="1000"
					bind:value={description}
					required
				></textarea>
			</label>

			<label class="form-control w-full">
				<div class="label"><span class="label-text font-medium">Your Character Origin</span></div>
				<input
					type="text"
					name="characterOrigin"
					placeholder="https://pdnode.com"
					class="input-bordered input w-full"
					autocomplete="name"
					bind:value={origin}
				/>
			</label>

			<div class="form-control w-full max-w-xs">
				<label class="label" for="file-upload">
					<span class="label-text">Your Character Pictures (Option)</span>
				</label>

				<input
					type="file"
					id="file-upload"
					class="file-input-bordered file-input w-full file-input-primary"
					accept="image/*"
					multiple
					bind:files={pictures}
				/>
			</div>

			<TagInput title="Your Character Tags (Option)" bind:tags />

			<p class="text-red-500">{errorText}</p>

			<div class="form-control mt-6">
				<button type="submit" class="btn w-full text-lg btn-primary" onclick={createCharacter}
					>Create</button
				>
			</div>
			<div class="modal-action">
				<form method="dialog">
					<!-- if there is a button in form, it will close the modal -->
					<button class="btn">Close</button>
				</form>
			</div>
		</div>
	</div>
</dialog>
