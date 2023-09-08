<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateS3ToBigquerySbLogTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('s3_to_bigquery_sb_logs', function (Blueprint $table) {
            $table->string('s3_key', 512)->primary();
            $table->enum('status', ['new', 'locked', 'finished'])->default('new');
            $table->boolean('error')->default(0);
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));

            $table->index('s3_key');
            $table->index(['status', 'error']);
            $table->index(['error', 'status']);
        });

        Schema::create('s3_to_bigquery_processes', function (Blueprint $table) {
            $table->string('process_uuid', 64)->primary();
            $table->enum('status', ['new', 'terminated', 'finished'])->default('new');
            $table->text('error')->nullable();
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));

            $table->index('process_uuid');
            $table->index('status');
        });

        Schema::create('s3_to_bigquery_exceeded', function (Blueprint $table) {
            $table->string('uuid', 64);
            $table->string('target', 512)->nullable();
            $table->string('key', 32);
            $table->bigInteger('timestamp')->default(0);
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));

            $table->primary(['uuid', 'target']);
            $table->index(['key', 'target', 'timestamp']);
            $table->index(['key', 'timestamp']);
        });

        Schema::create('s3_to_bigquery_load_files', function (Blueprint $table) {
            $table->string('table_id', 512)->primary();
            $table->string('process_uuid', 64)->nullable();
            $table->timestamp('updated_at')->default(DB::raw('CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP'));
            $table->timestamp('created_at')->default(DB::raw('CURRENT_TIMESTAMP'));
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('s3_to_bigquery_sb_logs');
        Schema::dropIfExists('s3_to_bigquery_processes');
        Schema::dropIfExists('s3_to_bigquery_exceeded');
        Schema::dropIfExists('s3_to_bigquery_load_files');
    }
}
